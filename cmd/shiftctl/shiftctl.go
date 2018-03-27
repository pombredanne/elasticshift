/*
Copyright 2018 The Elasticshift Authors.
*/
package main

import (
	"fmt"
	"os"

	"gitlab.com/conspico/elasticshift/pkg/shiftctl/cmd"
)

func main() {

	defaultCmd := cmd.NewDefaultCommand()
	if err := defaultCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
