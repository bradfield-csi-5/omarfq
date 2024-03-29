package leveldb

import (
	"bytes"
	"errors"

	"slices"

	"github.com/omarfq/leveldb/iterator"
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
}

type entry struct {
	key   []byte
	value []byte
}

type LevelDb struct {
	entries []entry
}

func NewLevelDb() *LevelDb {
	entries := make([]entry, 0)

	return &LevelDb{
		entries: entries,
	}
}

func (ldb *LevelDb) Get(key []byte) ([]byte, error) {
	idx, found := slices.BinarySearchFunc(ldb.entries, key, func(entry entry, target []byte) int {
		return bytes.Compare(entry.key, target)
	})

	if found {
		return ldb.entries[idx].value, nil
	}

	return nil, errors.New("Key not found")
}

func (ldb *LevelDb) Has(key []byte) (bool, error) {
	_, found := slices.BinarySearchFunc(ldb.entries, key, func(entry entry, target []byte) int {
		return bytes.Compare(entry.key, target)
	})

	return found, nil
}

func (ldb *LevelDb) Put(key, value []byte) error {
	idx, found := slices.BinarySearchFunc(ldb.entries, key, func(entry entry, target []byte) int {
		return bytes.Compare(entry.key, target)
	})

	newEntry := entry{key, value}

	if found {
		ldb.entries[idx].value = value
	} else {
		ldb.entries = append(ldb.entries, newEntry)
	}

	slices.SortFunc(ldb.entries, func(entry1, entry2 entry) int {
		return bytes.Compare(entry1.key, entry2.key)
	})

	return nil

}

func (ldb *LevelDb) Delete(key []byte) error {
	idx, found := slices.BinarySearchFunc(ldb.entries, key, func(entry entry, target []byte) int {
		return bytes.Compare(entry.key, target)
	})

	if found {
		ldb.entries = append(ldb.entries[0:idx], ldb.entries[idx+1:]...)
	}

	return errors.New("Could not delete entry. Provided key not found")
}

func (ldb *LevelDb) RangeScan(start, end []byte) (iterator.Iterator, error) {
	idxStart, foundStart := slices.BinarySearchFunc(ldb.entries, start, func(entry entry, target []byte) int {
		return bytes.Compare(entry.key, target)
	})

	idxEnd, foundEnd := slices.BinarySearchFunc(ldb.entries, end, func(entry entry, target []byte) int {
		return bytes.Compare(entry.key, target)
	})

	if !foundStart || !foundEnd {
		return nil, errors.New("One or both of the provided range keys could not be found")
	}

	newIterator := iterator.NewIter()
	rangeSlice := ldb.entries[idxStart : idxEnd+1]

	for _, entry := range rangeSlice {
		newTuple := iterator.Tuple{
			Key:   entry.key,
			Value: entry.value,
		}
		newIterator.Tuples = append(newIterator.Tuples, newTuple)
	}

	return newIterator, nil
}
