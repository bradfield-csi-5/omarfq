package leveldb

import (
	"errors"
	"fmt"

	"github.com/omarfq/leveldb/iterator"
	"github.com/omarfq/leveldb/skiplist"
)

type DB interface {
	// Get gets the value for the given key. It returns an error if the
	// DB does not contain the key.
	Get(key []byte) (value []byte, err error)

	// Has returns true if the DB contains the given key.
	Has(key []byte) (ret bool, err error)

	// Put sets the value for the given key. It overwrites any previous value
	// for that key; a DB is not a multi-map.
	Put(key, value []byte) error

	// Delete deletes the value for the given key.
	Delete(key []byte) error

	// RangeScan returns an Iterator (see below) for scanning through all
	// key-value pairs in the given range, ordered by key ascending.
	RangeScan(start, limit []byte) (iterator.Iterator, error)

	// Dump returns all the records in the database in a readable format
	Dump() (iterator.Iterator, error)
}

type LevelDb struct {
	entries *skiplist.SkipList
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
	ldb.entries.Insert(key, value)
	return nil
}

func (ldb *LevelDb) Delete(key []byte) error {
	err := ldb.entries.Delete(key)
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
