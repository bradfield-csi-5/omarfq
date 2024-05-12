package store

import (
	"encoding/binary"
	"fmt"
	pb "github.com/omarfq/kvstore/api/v1"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const PRIMARY_FILE_PATH = "data/kvstore.dat"
const SECONDARY_FILE_PATH = "data/kvstore_backup.dat"

type KVStore struct {
	PrimaryNode   *os.File
	SecondaryNode *os.File
}

func KVStoreInit() (*KVStore, error) {
	primaryDir := filepath.Dir(PRIMARY_FILE_PATH)
	if _, err := os.Stat(primaryDir); os.IsNotExist(err) {
		os.MkdirAll(primaryDir, 0752)
	}

	secondaryDir := filepath.Dir(SECONDARY_FILE_PATH)
	if _, err := os.Stat(secondaryDir); os.IsNotExist(err) {
		os.MkdirAll(secondaryDir, 0752)
	}

	primaryFile, err := os.OpenFile(PRIMARY_FILE_PATH, os.O_RDWR|os.O_CREATE, 0641)
	if err != nil {
		return nil, err
	}

	secondaryFile, err := os.OpenFile(SECONDARY_FILE_PATH, os.O_RDWR|os.O_CREATE, 0641)
	if err != nil {
		return nil, err
	}

	return &KVStore{
		PrimaryNode:   primaryFile,
		SecondaryNode: secondaryFile,
	}, nil
}

func (store *KVStore) Get(key *pb.Data) (string, error) {
	// Set file pointer to the beginning of the File
	_, err := store.PrimaryNode.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	prefixBuf := make([]byte, 8)

	for {
		// Read the FULL prefix length
		if _, err := io.ReadFull(store.PrimaryNode, prefixBuf); err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		// Decode the prefix length to an int
		recordLength := binary.BigEndian.Uint64(prefixBuf)

		// Read data
		dataBuf := make([]byte, recordLength)

		if _, err := io.ReadFull(store.PrimaryNode, dataBuf); err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		// Unmarshall binary data
		data := &pb.Data{}
		if err := proto.Unmarshal(dataBuf, data); err != nil {
			return "", fmt.Errorf("failed to unmarshal record data: %w", err)
		}

		// Compare keys
		if key.Key == data.Key {
			return data.Value, nil
		}
	}

	return "", fmt.Errorf("the key: %s does not exist in the key-value store", key.Key)
}

func (store *KVStore) Set(keyValue *pb.Data) error {
	// Check if key exists, if it does update it
	existingValue, err := store.Get(keyValue)
	if err == nil {
		if existingValue == keyValue.Value {
			return nil // No update required
		}
	} else if !strings.Contains(err.Error(), "does not exist") {
		return fmt.Errorf("failed to check existing key: %w", err)
	}

	// Marshal the new key-value pair
	dataBytes, err := proto.Marshal(keyValue)
	if err != nil {
		return fmt.Errorf("failed to serialize key-value pair: %w", err)
	}

	lengthBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(lengthBytes, uint64(len(dataBytes)))

	err = appendToEndOfFile(lengthBytes, dataBytes, store.PrimaryNode)
	if err != nil {
		return err
	}

	err = appendToEndOfFile(lengthBytes, dataBytes, store.SecondaryNode)
	if err != nil {
		return err
	}

	return nil
}

func appendToEndOfFile(lengthBytes, databytes []byte, node *os.File) error {
	// Append to the end of the file
	_, err := node.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("failed to seek to end of file: %w", err)
	}

	if _, err := node.Write(lengthBytes); err != nil {
		return fmt.Errorf("failed to write length prefix: %w", err)
	}
	if _, err := node.Write(databytes); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	// Flush to disk
	if err := node.Sync(); err != nil {
		return fmt.Errorf("failed to flush data to disk: %w", err)
	}
	return nil
}
