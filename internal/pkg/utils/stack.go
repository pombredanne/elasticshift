/*
Copyright 2018 The Elasticshift Authors.
*/
package utils

import "fmt"

type Stack interface {
	Push(item string)
	Pop() string
	Last() string
	Print()
}

type stack struct {
	items []string
	size  int
}

func NewStack() Stack {

	s := &stack{}
	s.items = []string{}
	return s
}

func (s *stack) Push(item string) {
	s.items = append(s.items, item)
	s.size = len(s.items)
}

func (s *stack) Last() string {

	if s.size == 0 {
		return ""
	}

	return s.items[s.size-1]
}

func (s *stack) Pop() string {

	if s.size == 0 {
		return ""
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
