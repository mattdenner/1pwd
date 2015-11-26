package opvault

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Band map[string]*Item

func parseBand(data []byte) (Band, error) {
	var (
		idx  int
		band Band
	)

	idx = bytes.IndexByte(data, '{')
	if idx < 0 {
		return nil, errors.New("invalid band data")
	}
	data = data[idx:]

	idx = bytes.LastIndexByte(data, '}')
	if idx < 0 {
		return nil, errors.New("invalid band data")
	}
	data = data[:idx+1]

	err := json.Unmarshal(data, &band)
	if err != nil {
		return nil, err
	}

	return band, nil
}

func (b Band) decryptOverView(p *Profile) error {
	for _, item := range b {
		err := item.decryptOverView(p)
		if err != nil {
			return err
		}
	}
	return nil
}
