package iterator

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

// Assuming SkipListNode is defined in your skip list implementation
type SkipListNode struct {
	Key, Value []byte
	Forwards   []*SkipListNode // Array of forward
}

type SkipListIter struct {
	current *SkipListNode
}

// NewSkipListIter creates a new iterator for a skip list starting at the given node
func NewSkipListIter(startNode *SkipListNode) *SkipListIter {
	return &SkipListIter{current: startNode}
}

// Next moves the iterator to the next node
func (iter *SkipListIter) Next() bool {
	if iter.current != nil && iter.current.Forwards[0] != nil {
		iter.current = iter.current.Forwards[0]
		return true
	}
	return false
}

// Key returns the key of the current node
func (iter *SkipListIter) Key() []byte {
	if iter.current != nil {
		return iter.current.Key
	}
	return nil
}

// Value returns the value of the current node
func (iter *SkipListIter) Value() []byte {
	if iter.current != nil {
		return iter.current.Value
	}
	return nil
}

// Error could potentially be implemented to handle any errors encountered during iteration
func (iter *SkipListIter) Error() error {
	// Implementation depends on if/how you want to handle errors
	return nil
}
