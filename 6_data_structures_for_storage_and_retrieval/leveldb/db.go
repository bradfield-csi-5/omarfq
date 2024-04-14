package leveldb

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/omarfq/leveldb/iterator"
	"github.com/omarfq/leveldb/skiplist"
	"github.com/omarfq/leveldb/wal"
)

type LevelDb struct {
	entries *skiplist.SkipList
	wal     *wal.WAL
}

func NewLevelDb(key, val []byte) *LevelDb {
	// Create a new SkipList instance
	newSkipList := skiplist.NewSkipList()

	// Create the new LevelDb instance with the newly created SkipList
	ldb := &LevelDb{
		entries: newSkipList,
	}

	// If an initial key and value are provided, insert them into the skip list
	if key != nil && val != nil {
		ldb.entries.Insert(key, val)
	}

	return ldb
}

func (ldb *LevelDb) Get(key []byte) ([]byte, error) {
	node := ldb.entries.Search(key)
	if node == nil {
		return nil, errors.New("Key not found")
	}

	return node.Value, nil
}

func (ldb *LevelDb) Has(key []byte) (bool, error) {
	node := ldb.entries.Search(key)
	if node == nil {
		return false, nil
	}
	return true, nil
}

func (ldb *LevelDb) Put(key, value []byte) error {
	err := ldb.wal.Write(wal.OpPut, key, value)
	if err != nil {
		return err
	}
	ldb.entries.Insert(key, value)
	return nil
}

func (ldb *LevelDb) Delete(key []byte) error {
	err := ldb.wal.Write(wal.OpDelete, key, nil)
	if err != nil {
		return err
	}
	err = ldb.entries.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (ldb *LevelDb) RangeScan(start, end []byte) (iterator.Iterator, error) {
	startNode := ldb.entries.Search(start)
	if startNode == nil {
		return nil, fmt.Errorf("Start key not found")
	}

	return iterator.NewSkipListIterator(startNode), nil
}

func (ldb *LevelDb) RecoverFromWAL() error {
	file, err := os.Open("wal.go")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		op, err := reader.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		keyLenBuf := make([]byte, 4) // Key length pos in binary format is 4 bytes max
		_, err = io.ReadFull(reader, keyLenBuf)
		if err != nil {
			return err
		}

		keyLen := binary.BigEndian.Uint32(keyLenBuf)

		var valueLen uint32
		if op == wal.OpPut {
			valueLenBuf := make([]byte, 4)
			_, err = io.ReadFull(reader, valueLenBuf)
			if err != nil {
				return err
			}
			valueLen = binary.BigEndian.Uint32(valueLenBuf)
		}

		key := make([]byte, keyLen)
		_, err = io.ReadFull(reader, key)
		if err != nil {
			return err
		}

		var value []byte
		if op == wal.OpPut {
			value = make([]byte, valueLen)
			_, err := io.ReadFull(reader, value)
			if err != nil {
				return err
			}
		}

		// Decide which op to execute
		switch op {
		case wal.OpPut:
			ldb.Put(key, value)
		case wal.OpDelete:
			ldb.Delete(key)
		}
	}
	return nil
}
