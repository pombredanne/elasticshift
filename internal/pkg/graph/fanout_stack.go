/*
Copyright 2018 The Elasticshift Authors.
*/
package graph

import "fmt"

type FanNStack interface {
	Push(item *FanN)
	Pop() *FanN
	Last() *FanN
	Print()
	Size() int
}

type stack struct {
	items []*FanN
	size  int
}

// NewFanNStack ..
func NewFanNStack() FanNStack {

	s := &stack{}
	s.items = []*FanN{}
	return s
}

func (s *stack) Push(item *FanN) {
	s.items = append(s.items, item)
	s.size = len(s.items)
}

func (s *stack) Last() *FanN {

	if s.size == 0 {
		return nil
	}

	return s.items[s.size-1]
}

func (s *stack) Size() int {
	return s.size
}

func (s *stack) Pop() *FanN {

	if s.size == 0 {
		return nil
	}

	lastIndex := s.size - 1
	item := s.items[lastIndex]
	s.items = s.items[0:lastIndex]
	s.size = lastIndex

	return item
}

func (s *stack) Print() {
	fmt.Println(s.items)
}
