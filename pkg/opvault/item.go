package opvault

import (
	"encoding/json"
	"net/url"
	"sort"
)

type Item struct {
	UUID     string
	Category Category
	Fave     int
	Trashed  bool
	Folder   string
	O        []byte
	K        []byte
	D        []byte
	HMAC     []byte
	Tx       int64
	Updated  int64
	Created  int64

	Data *struct {
		UUID     string   `json:"uuid,omitempty"`
		Category Category `json:"category,omitempty"`

		// Overview
		Title  string `json:"title,omitempty"`
		URL    string `json:"url,omitempty"`
		Domain string `json:"domain,omitempty"`
		URLs   []struct {
			U string `json:"url,omitempty"`
			L string `json:"label,omitempty"`
		}

		// Data
		BackupKeys [][]byte `json:"backupKeys"`
		Password   string   `json:"password,omitempty"`
		Fields     []struct {
			Type        string `json:"type,omitempty"`
			Name        string `json:"name,omitempty"`
			Designation string `json:"designation,omitempty"`
			Value       string `json:"value,omitempty"`
		} `json:"fields,omitempty"`
	}
}

func (i *Item) decryptOverView(p *Profile) error {
	dst, err := decrypt(nil, i.O, p.overviewEncKey, p.overviewMacKey)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dst, &i.Data)
	if err != nil {
		return err
	}

	i.Data.UUID = i.UUID
	i.Data.Category = i.Category

	if i.Data.URL != "" {
		u, err := url.Parse(i.Data.URL)
		if err == nil {
			i.Data.Domain = u.Host
		}
	}

	return nil
}

func (i *Item) Decrypt(v *Vault) error {
	return i.decryptData(v.profile)
}

func (i *Item) decryptData(p *Profile) error {
	dstKey, err := decryptKey(nil, i.K, p.masterEncKey, p.masterMacKey)
	if err != nil {
		return err
	}

	dst, err := decrypt(nil, i.D, dstKey[:32], dstKey[32:])
	if err != nil {
		return err
	}

	err = json.Unmarshal(dst, &i.Data)
	if err != nil {
		return err
	}

	// var buf bytes.Buffer
	// json.Indent(&buf, dst, "", "  ")
	// fmt.Fprintf(os.Stderr, "data: %s\n", buf.String())

	return nil
}

func (i *Item) Extract(field string) (string, bool) {
	var (
		v string
		f bool
	)

	switch field {

	case "url":
		if !f && i.Data.URL != "" {
			v, f = i.Data.URL, true
		}
		if !f && len(i.Data.URLs) > 0 {
			v, f = i.Data.URLs[0].U, true
		}
		return v, f

	case "username":
		if !f {
			v, f = i.extractFieldByDesignation("username")
		}
		if !f {
			v, f = i.extractFieldByName("username")
		}
		if !f {
			v, f = i.extractFieldByName("login")
		}
		return v, f

	case "password":
		if !f && i.Data.Password != "" {
			v, f = i.Data.Password, true
		}
		if !f {
			v, f = i.extractFieldByDesignation("password")
		}
		if !f {
			v, f = i.extractFieldByName("password")
		}
		return v, f

	default:
		if !f {
			v, f = i.extractFieldByDesignation(field)
		}
		if !f {
			v, f = i.extractFieldByName(field)
		}
		return v, f

	}
}

func (i *Item) extractFieldByName(field string) (string, bool) {
	for _, f := range i.Data.Fields {
		if f.Name == field {
			return f.Value, true
		}
	}

	return "", false
}

func (i *Item) extractFieldByDesignation(field string) (string, bool) {
	for _, f := range i.Data.Fields {
		if f.Designation == field {
			return f.Value, true
		}
	}

	return "", false
}

type byTitle []*Item

func (s byTitle) Len() int           { return len(s) }
func (s byTitle) Less(i, j int) bool { return s[i].Data.Title < s[j].Data.Title }
func (s byTitle) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type byDomain []*Item

func (s byDomain) Len() int           { return len(s) }
func (s byDomain) Less(i, j int) bool { return s[i].Data.Domain < s[j].Data.Domain }
func (s byDomain) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type multiSort []sort.Interface

func (s multiSort) Len() int      { return s[0].Len() }
func (s multiSort) Swap(i, j int) { s[0].Swap(i, j) }
func (s multiSort) Less(i, j int) bool {
	for _, iface := range s {
		if iface.Less(i, j) {
			return true
		} else if iface.Less(j, i) {
			return false
		}
	}
	return false
}
