package opvault

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

type Profile struct {
	UUID          string
	UpdatedAt     int64
	CreatedAt     int64
	LastUpdatedBy string
	ProfileName   string
	PasswordHint  string
	Iterations    int
	Salt          []byte
	OverviewKey   []byte
	MasterKey     []byte

	masterEncKey   []byte
	masterMacKey   []byte
	overviewEncKey []byte
	overviewMacKey []byte
}

func parseProfile(data []byte) (*Profile, error) {
	var (
		idx     int
		profile *Profile
	)

	idx = bytes.IndexByte(data, '{')
	if idx < 0 {
		return nil, errors.New("invalid profile data")
	}
	data = data[idx:]

	idx = bytes.LastIndexByte(data, '}')
	if idx < 0 {
		return nil, errors.New("invalid profile data")
	}
	data = data[:idx+1]

	err := json.Unmarshal(data, &profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (p *Profile) setMasterPassword(pwd string) error {

	var (
		dk            = pbkdf2.Key([]byte(pwd), p.Salt, p.Iterations, 64, sha512.New)
		derivedEncKey = dk[:32]
		derivedMacKey = dk[32:]
	)

	masterKey, err := decrypt(nil, p.MasterKey, derivedEncKey, derivedMacKey)
	if err != nil {
		return err
	}

	overviewKey, err := decrypt(nil, p.OverviewKey, derivedEncKey, derivedMacKey)
	if err != nil {
		return err
	}

	mac := sha512.New()
	mac.Write(masterKey)
	macData := mac.Sum(nil)
	p.masterEncKey = macData[:32]
	p.masterMacKey = macData[32:]

	mac.Reset()
	mac.Write(overviewKey)
	macData = mac.Sum(nil)
	p.overviewEncKey = macData[:32]
	p.overviewMacKey = macData[32:]

	return nil
}
