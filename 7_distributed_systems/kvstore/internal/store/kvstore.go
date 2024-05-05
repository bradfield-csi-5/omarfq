package store

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileKVStore struct {
	path string
}

func NewFileKVStore(path string) *FileKVStore {
	return &FileKVStore{
		path: path,
	}
}

func (store *FileKVStore) Get(key string) (string, error) {
	file, err := os.Open(store.path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data map[string]string
	if err := decoder.Decode(&data); err != nil {
		return "", err
	}

	val, ok := data[key]
	if !ok {
		return "", fmt.Errorf("Key %q not found in the file", key)
	}

	return val, nil
}

func (store *FileKVStore) Set(key, value string) error {
	dir := filepath.Dir(store.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	file, err := os.OpenFile(store.path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data map[string]string
	if err := decoder.Decode(&data); err != nil {
		if err != io.EOF {
			return err
		}
		data = make(map[string]string)
	}

	data[key] = value

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}
