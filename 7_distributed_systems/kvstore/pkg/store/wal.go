package store

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
)

const (
	WAL_PATH  = "data/kvstore.log"
	OP_SET    = 0x01
	OP_DELETE = 0x02
)

type WAL struct {
	file *os.File
}

func NewWAL() (*WAL, error) {
	dir := filepath.Dir(WAL_PATH)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	f, err := os.OpenFile(WAL_PATH, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{
		file: f,
	}, nil
}

func (wal *WAL) Write(op byte, key, value []byte) error {
	buf := new(bytes.Buffer)

	if err := buf.WriteByte(op); err != nil {
		if err != nil {
			return err
		}
		return err
	}

	if _, err := buf.Write(encodeUint32(uint32(len(key)))); err != nil {
		return err
	}
	if op == OP_SET {
		if _, err := buf.Write(encodeUint32(uint32(len(value)))); err != nil {
			return err
		}
	}

	if _, err := buf.Write(key); err != nil {
		return err
	}
	if op == OP_SET {
		if _, err := buf.Write(value); err != nil {
			return err
		}
	}

	if _, err := wal.file.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

// TODO: Implement GetNext with iterator pattern for
// recovering a node from the WAL
func (wal *WAL) GetNext() {}

func encodeUint32(u uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, u)
	return buf
}
