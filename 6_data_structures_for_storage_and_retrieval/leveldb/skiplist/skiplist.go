package skiplist

import (
	"bytes"
	"math/rand"
	"time"
)

const (
	MaxLevel    int     = 16
	Probability float64 = 0.5
)

type SkipListNode struct {
	Key, Value []byte
	Forwards   []*SkipListNode
}

type SkipList struct {
	Head     *SkipListNode
	Level    int
	MaxLevel int
	rnd      *rand.Rand
}

func NewSkipList() *SkipList {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	head := &SkipListNode{
		Key:      nil,
		Value:    nil,
		Forwards: make([]*SkipListNode, MaxLevel),
	}

	return &SkipList{
		Head:     head,
		Level:    0,
		MaxLevel: MaxLevel,
		rnd:      rnd,
	}
}

func (sl *SkipList) randomLevel() int {
	level := 1
	for sl.rnd.Float64() < Probability && level < sl.MaxLevel {
		level++
	}
	return level
}

func (sl *SkipList) Insert(key, value []byte) {
	update := make([]*SkipListNode, sl.MaxLevel)
	current := sl.Head

	for i := sl.Level - 1; i >= 0; i-- {
		for current.Forwards[i] != nil && bytes.Compare(current.Forwards[i].Key, key) < 0 {
			current = current.Forwards[i]
		}
		update[i] = current
	}

	level := sl.randomLevel()

	if level > sl.Level {
		for i := sl.Level; i < level; i++ {
			update[i] = sl.Head
		}
		sl.Level = level
	}

	newNode := &SkipListNode{
		Key:      key,
		Value:    value,
		Forwards: make([]*SkipListNode, level),
	}

	for i := 0; i < level; i++ {
		newNode.Forwards[i] = update[i].Forwards[i]
		update[i].Forwards[i] = newNode
	}
}

func (sl *SkipList) Delete(key []byte) {
	update := make([]*SkipListNode, sl.MaxLevel)
	current := sl.Head

	for i := sl.Level - 1; i >= 0; i-- {
		for current.Forwards[i] != nil && bytes.Compare(current.Forwards[i].Key, key) < 0 {
			current = current.Forwards[i]
		}
		update[i] = current
	}

	current = current.Forwards[0]

	if current != nil && bytes.Compare(current.Key, key) == 0 {
		for i := 0; i < sl.Level; i++ {
			if update[i].Forwards[i] != current {
				break
			}
			update[i].Forwards[i] = current.Forwards[i]
		}

		for sl.Level > 1 && sl.Head.Forwards[sl.Level-1] == nil {
			sl.Level--
		}
	}
}

func (sl *SkipList) Update(key, value []byte) {
	current := sl.Head

	for i := sl.Level - 1; i >= 0; i-- {
		for current.Forwards[i] != nil && bytes.Compare(current.Forwards[i].Key, key) < 0 {
			current = current.Forwards[i]
		}
	}

	current = current.Forwards[0]

	if current != nil && bytes.Compare(current.Key, key) == 0 {
		current.Value = value
	} else {
		sl.Insert(key, value)
	}
}
