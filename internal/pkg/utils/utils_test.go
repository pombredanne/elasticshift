/*
Copyright 2018 The Elasticshift Authors.
*/
package utils

import (
	"fmt"
	"testing"
)

func TestPathExist(t *testing.T) {

	exist, err := PathExist("/Users/ghazni/Brewfile")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(exist)
}
