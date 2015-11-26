package opvault

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Folders map[string]*Folder

type Folder struct {
	UUID     string
	Parent   string
	Updated  int64
	Created  int64
	Tx       int64
	Smart    bool
	Overview []byte
}

func parseFolders(data []byte) (Folders, error) {
	var (
		idx     int
		folders Folders
	)

	idx = bytes.IndexByte(data, '{')
	if idx < 0 {
		return nil, errors.New("invalid folders data")
	}
	data = data[idx:]

	idx = bytes.LastIndexByte(data, '}')
	if idx < 0 {
		return nil, errors.New("invalid folders data")
	}
	data = data[:idx+1]

	err := json.Unmarshal(data, &folders)
	if err != nil {
		return nil, err
	}

	return folders, nil
}
