package lox

import (
	"iter"
	"slices"
)

type Stack[I any] struct {
	items []I
}

func (s *Stack[I]) len() int {
	return len(s.items)
}

func (s *Stack[I]) isEmpty() bool {
	return s.len() == 0
}

func (s *Stack[I]) push(item I) {
	s.items = append(s.items, item)
}

func (s *Stack[I]) pop() I {
	result := s.items[s.len()-1]
	s.items = s.items[:s.len()-1]
	return result
}

func (s *Stack[I]) peek() I {
	return s.items[s.len()-1]
}

func (s *Stack[I]) walk() iter.Seq2[int, I] {
	return slices.Backward(s.items)
}
