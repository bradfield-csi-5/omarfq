package executor

import (
	"iter/Iterator"
)

type SelectionNode struct {
	source    Iterator
	predicate func([]ColumnValue) bool
}

func (s *SelectionNode) Next() []ColumnValue {
	for {
		candidate := s.source.Next()
		if s.predicate(candidate) {
			return candidate
		}
	}
}
