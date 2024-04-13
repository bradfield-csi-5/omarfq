package iterator

import (
	"errors"

	"github.com/omarfq/leveldb/skiplist"
)

type Iterator interface {
	// Next moves the iterator to the next key/value pair.
	// It returns false if the iterator is exhausted.
	Next() bool

	// Error returns any accumulated error. Exhausting all the key/value pairs
	// is not considered to be an error.
	Error() error

	// Key returns the key of the current key/value pair, or nil if done.
	Key() []byte

	// Value returns the value of the current key/value pair, or nil if done.
	Value() []byte
}

type SkipListIter struct {
	current *skiplist.SkipListNode
	err     error
}

func NewSkipListIterator(startNode *skiplist.SkipListNode) *SkipListIter {
	return &SkipListIter{current: startNode}
}

func (iter *SkipListIter) Next() bool {
	if iter.current != nil && iter.current.Forwards[0] != nil {
		iter.current = iter.current.Forwards[0]
		return true
	}
	if iter.current == nil {
		iter.err = errors.New("Iterator has reached end")
	}
	return false
}

func (iter *SkipListIter) Key() []byte {
	if iter.current != nil {
		return iter.current.Key
	}
	return nil
}

func (iter *SkipListIter) Value() []byte {
	if iter.current != nil {
		return iter.current.Value
	}
	return nil
}

func (iter *SkipListIter) Error() error {
	return iter.err
}
