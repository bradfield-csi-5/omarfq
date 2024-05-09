package store

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	pb "github.com/omarfq/kvstore/api/v1"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"path/filepath"
)

const FILE_PATH = "data/kvstore.dat"

type FileKVStore struct {
	File *os.File
}

func FileKVStoreInstance() (*FileKVStore, error) {
	dir := filepath.Dir(FILE_PATH)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	file, err := os.OpenFile(FILE_PATH, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &FileKVStore{
		File: file,
	}, nil
}

func (store *FileKVStore) Get(key *pb.Data) (string, error) {
	// Set file pointer to the beginning of the File
	_, err := store.File.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	prefixBuf := make([]byte, 8)

	for {
		// Read the FULL prefix length
		if _, err := io.ReadFull(store.File, prefixBuf); err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		// Decode the prefix length to an int
		recordLength := binary.BigEndian.Uint64(prefixBuf)

		// Read data
		dataBuf := make([]byte, recordLength)

		if _, err := io.ReadFull(store.File, dataBuf); err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		// Unmarshall binary data
		data := &pb.Data{}
		if err := proto.Unmarshal(dataBuf, data); err != nil {
			return "", fmt.Errorf("Failed to unmarshal record data: %w", err)
		}

		// Compare keys
		if key.Key == data.Key {
			return data.Value, nil
		}
	}

	return "", fmt.Errorf("The key: %s does not exist in the key-value store", key.Key)
}

func (store *FileKVStore) Set(kvpair *pb.Data) error {
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
