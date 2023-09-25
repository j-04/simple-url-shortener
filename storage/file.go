package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type internalStore struct {
	Version string            `json:"version"`
	Items   map[string]string `json:"items"`
}

type FileStore struct {
	filename string
}

func NewFileStore(filename string) (FileStore, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		is := internalStore{Version: "v1", Items: make(map[string]string)}
		raw, err := json.Marshal(is)
		if err != nil {
			return FileStore{}, fmt.Errorf("unable to generate json representation for file")
		}

		err = os.WriteFile(filename, raw, 0644)
		if err != nil {
			return FileStore{}, fmt.Errorf("Unable to persist file")
		}
	}
	return FileStore{filename: filename}, nil
}

func (f *FileStore) Add(shortenedUrl, longUrl string) error {
	raw, err := os.ReadFile(f.filename)
	if err != nil {
		return err
	}

	var is internalStore

	err = json.Unmarshal(raw, &is)
	if err != nil {
		return fmt.Errorf("Unable to parse incoming json store data. Err: %v", err)
	}

	_, ok := is.Items[shortenedUrl]
	if ok {
		return fmt.Errorf("Shortened url already stored")
	}
	is.Items[shortenedUrl] = longUrl
	modRaw, err := json.Marshal(is)
	if err != nil {
		return fmt.Errorf("unable to convert data to json representation")
	}

	err = os.WriteFile(f.filename, modRaw, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileStore) Remove(shortenedUrl string) error {
	raw, err := os.ReadFile(f.filename)
	if err != nil {
		return err
	}

	var is internalStore
	err = json.Unmarshal(raw, &is)
	if err != nil {
		return fmt.Errorf("Unalbe to parse incoming json store data. Err: %v", err)
	}

	delete(is.Items, shortenedUrl)
	modRaw, err := json.Marshal(is)
	if err != nil {
		return fmt.Errorf("Unable to convert data to json representation")
	}

	err = os.WriteFile(f.filename, modRaw, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileStore) Get(shortenedUrl string) (string, error) {
	raw, err := os.ReadFile(f.filename)
	if err != nil {
		return "", err
	}

	var is internalStore
	err = json.Unmarshal(raw, &is)
	if err != nil {
		return "", fmt.Errorf("Unable to parse incoming json store data. Err: %v", err)
	}

	longUrl, ok := is.Items[shortenedUrl]
	if !ok {
		return "", fmt.Errorf("No url available for that shortened url")
	}

	return longUrl, nil
}
