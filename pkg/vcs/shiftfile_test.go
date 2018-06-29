/*
Copyright 2018 The Elasticshift Authors.
*/
package vcs

import (
	"fmt"
	"testing"
)

func TestShiftfileFromGithub(t *testing.T) {

	url := "https://github.com/nshahm/hybrid.test.runner.git"
	branch := "master"
	data, err := GetShiftFile(GITHUB_DOT_COM, url, branch)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(data[:]))
}
