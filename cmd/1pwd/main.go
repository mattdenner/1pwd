package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"text/tabwriter"

	"github.com/bgentry/speakeasy"
	"github.com/fd/1pwd/pkg/opvault"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		id            string
		vaultPath     string
		extract       string
		query         string
		typeFilter    string
		coprocessType string
		coprocessCode string
		jsonFormat    bool
	)

	app := kingpin.New("1pwd", "A command-line tool for 1Password.").
		Author("Simon Menke").
		Version("1.0.0")
	app.Flag("vault", "Vault to read").Short('V').StringVar(&vaultPath)

	get := app.Command("get", "Get an entry")
	get.Arg("id", "ID of item.").Required().StringVar(&id)
	get.Arg("extract", "Field to extract").StringVar(&extract)
	get.Flag("json", "Print JSON formatted data").Short('j').BoolVar(&jsonFormat)

	search := app.Command("search", "Search for an entry")
	search.Arg("extract", "Field to extract").StringVar(&extract)
	search.Flag("type", "Entry type").Short('t').Default("login").EnumVar(&typeFilter,
		"any",
		opvault.LoginItem.TypeString(),
		opvault.CreditCardItem.TypeString(),
		opvault.SecureNoteItem.TypeString(),
		opvault.IdentityItem.TypeString(),
		opvault.PasswordItem.TypeString(),
		opvault.TombstoneItem.TypeString(),
		opvault.SoftwareLicenseItem.TypeString(),
		opvault.BankAccountItem.TypeString(),
		opvault.DatabaseItem.TypeString(),
		opvault.DriverLicenseItem.TypeString(),
		opvault.OutdoorLicenseItem.TypeString(),
		opvault.MembershipItem.TypeString(),
		opvault.PassportItem.TypeString(),
		opvault.RewardsItem.TypeString(),
		opvault.SSNItem.TypeString(),
		opvault.RouterItem.TypeString(),
		opvault.ServerItem.TypeString(),
		opvault.EmailItem.TypeString())
	search.Flag("query", "Initial query").Short('q').StringVar(&query)
	search.Flag("json", "Print JSON formatted data").Short('j').BoolVar(&jsonFormat)
	search.Flag("coprocess", "Coprocess code").Hidden().StringVar(&coprocessCode)

	coprocess := app.Command("coprocess", "iTerm coprocess").Hidden()
	coprocess.Arg("type", "coprocess type").Default("iterm").EnumVar(&coprocessType, "iterm")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case get.FullCommand():
		doGet(openVault(vaultPath), id, extract, "", jsonFormat)
	case search.FullCommand():
		doSearch(openVault(vaultPath), query, typeFilter, extract, coprocessCode, jsonFormat)
	case coprocess.FullCommand():
		doCoprocess(coprocessType)
	}
}

func openVault(vaultPath string) *opvault.Vault {
	if vaultPath == "" {
		vaults, err := opvault.LookupVaults()
		assert(err)

		if len(vaults) == 0 {
			abortf("no vaults found")
		}

		vaultPath = vaults[0]
	}

	pwd, err := speakeasy.FAsk(os.Stderr, "Master Password: ")
	assert(err)

	vault, err := opvault.Open(vaultPath, pwd)
	assert(err)

	return vault
}

func doSearch(vault *opvault.Vault, query, typeFilter, extract, coprocessCode string, jsonFormat bool) {
	if typeFilter == "any" {
		typeFilter = ""
	}
	cat := opvault.FromTypeString(typeFilter)

	results := vault.All()

	var bufIn bytes.Buffer
	var bufOut bytes.Buffer

	tabw := tabwriter.NewWriter(&bufIn, 8, 8, 2, ' ', tabwriter.StripEscape)
	for _, result := range results {
		if result.Trashed {
			continue
		}
		if typeFilter != "" && result.Category != cat {
			continue
		}

		fmt.Fprintf(tabw,
			field("%s", "")+
				field("%s", "2")+
				field("%s", "")+
				field("%s", "blue")+
				field("%s", "yellow")+
				"\n",
			result.UUID,
			result.UUID[:8],
			result.Category.String(),
			trunc(result.Data.Domain, 32),
			result.Data.Title,
		)
	}
	tabw.Flush()

	cmd := exec.Command("fzf", "--ansi", "--with-nth=2..", "--nth=4..,3,1", "--query="+query)
	cmd.Env = os.Environ()
	cmd.Stdin = &bufIn
	cmd.Stdout = &bufOut
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	assert(err)

	var id string
	fmt.Fscan(&bufOut, &id)
	if id != "" {
		doGet(vault, id, extract, coprocessCode, jsonFormat)
	}
}

func doGet(vault *opvault.Vault, id, extract, coprocessCode string, jsonFormat bool) {
	if coprocessCode != "" {
		extract = "password"
		jsonFormat = false
	}

	item, err := vault.Get(id)
	assert(err)

	err = item.Decrypt(vault)
	assert(err)

	var (
		v interface{} = item.Data
		f             = true
	)

	if extract != "" {
		v, f = item.Extract(extract)
		if !f {
			abortf("field %q not found", extract)
		}
	}

	if coprocessCode != "" {
		conn, err := net.Dial("unix", "/tmp/"+coprocessCode+".sock")
		assert(err)
		defer conn.Close()

		var data = struct {
			Password string
		}{
			Password: v.(string),
		}

		err = json.NewEncoder(conn).Encode(&data)
		assert(err)

	} else if jsonFormat {
		err = json.NewEncoder(os.Stdout).Encode(v)
		assert(err)
	} else {
		if extract != "" {
			fmt.Printf("%s\n", v)
		} else {
			fmt.Printf("id:       %s\n", item.UUID)
			if v, f := item.Extract("url"); f {
				fmt.Printf("url:      %s\n", v)
			}
			if v, f := item.Extract("username"); f {
				fmt.Printf("username: %s\n", v)
			}
			if v, f := item.Extract("password"); f {
				fmt.Printf("password: %s\n", v)
			}
		}
	}
}

func trunc(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}

	return s
}

func field(format, color string) string {
	if color == "" {
		return format + "\t"
	}

	switch color {
	case "black":
		color = "30"
	case "red":
		color = "31"
	case "green":
		color = "32"
	case "yellow":
		color = "33"
	case "blue":
		color = "34"
	case "magenta":
		color = "35"
	case "cyan":
		color = "36"
	case "white":
		color = "37"
	}

	return "\x1B[" + color + "m" + format + "\x1B[0m\t"
}

func abortf(format string, args ...interface{}) {
	assert(fmt.Errorf(format, args...))
}

func assert(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "1pwd: error: %s\n", err)
		os.Exit(1)
	}
}
