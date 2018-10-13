/*
Copyright 2018 The Elasticshift Authors.
*/
package utils

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {

	stk := NewStack()
	stk.Push("1")
	stk.Push("1:1")
	stk.Push("1:1:1")

	stk.Print()

	var item string
	item = stk.Pop()

	fmt.Println("Item=", item)
	stk.Print()

	item = stk.Pop()

	fmt.Println("Item=", item)
	stk.Print()

	item = stk.Pop()

	fmt.Println("Item=", item)
	stk.Print()

	item = stk.Pop()

	fmt.Println("Item=", item)
	stk.Print()
}
