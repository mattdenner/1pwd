package opvault

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
)

type Vault struct {
	profile *Profile
	folders Folders
	bands   [16]Band
}

func Open(path, master string) (*Vault, error) {
	var (
		vault = &Vault{}
		data  []byte
		err   error
	)

	data, err = ioutil.ReadFile(filepath.Join(path, "default", "profile.js"))
	if err != nil {
		return nil, err
	}

	vault.profile, err = parseProfile(data)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadFile(filepath.Join(path, "default", "folders.js"))
	if err != nil {
		return nil, err
	}

	vault.folders, err = parseFolders(data)
	if err != nil {
		return nil, err
	}

	for i := '0'; i <= '9'; i++ {
		data, err = ioutil.ReadFile(filepath.Join(path, "default", fmt.Sprintf("band_%c.js", i)))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		band, err := parseBand(data)
		if err != nil {
			return nil, err
		}

		vault.bands[i-'0'] = band
	}

	for i := 'A'; i <= 'F'; i++ {
		data, err = ioutil.ReadFile(filepath.Join(path, "default", fmt.Sprintf("band_%c.js", i)))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		band, err := parseBand(data)
		if err != nil {
			return nil, err
		}

		vault.bands[(i-'A')+10] = band
	}

	err = vault.profile.setMasterPassword(master)
	if err != nil {
		return nil, err
	}

	err = vault.decryptOverView()
	if err != nil {
		return nil, err
	}

	return vault, nil
}

func (v *Vault) Get(itemID string) (*Item, error) {
	if itemID == "" {
		return nil, os.ErrNotExist
	}

	bandID, err := strconv.ParseInt(itemID[:1], 16, 8)
	if err != nil {
		return nil, err
	}
	if bandID < 0 || bandID >= 16 {
		return nil, os.ErrNotExist
	}

	band := v.bands[bandID]
	if band == nil {
		return nil, os.ErrNotExist
	}

	item := band[itemID]
	if item == nil {
		return nil, os.ErrNotExist
	}

	return item, nil
}

func (v *Vault) All() []*Item {
	var results = make([]*Item, 0, 4096)

	for _, band := range v.bands {
		for _, item := range band {
			results = append(results, item)
		}
	}

	sort.Sort(multiSort{byDomain(results), byTitle(results)})

	return results
}

func (v *Vault) decryptOverView() error {
	for _, band := range v.bands {
		err := band.decryptOverView(v.profile)
		if err != nil {
			return err
		}
	}
	return nil
}

func LookupVaults() ([]string, error) {
	var home string

	if home == "" {
		home = os.Getenv("HOME")
	}

	if home == "" {
		u, err := user.Current()
		if err != nil {
			return nil, err
		}
		home = u.HomeDir
	}

	entries, err := filepath.Glob(filepath.Join(home, "Dropbox*", "*.opvault"))
	if err != nil {
		return nil, err
	}

	return entries, nil
}
