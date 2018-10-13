/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"fmt"
	"log"

	"gitlab.com/conspico/elasticshift/internal/pkg/graph"
)

var (
	SHELL = "shell"
)

func (b *builder) invokePlugin(n *graph.N) (string, error) {

	if graph.START == n.Name || graph.END == n.Name {
		return "", nil
	}

	var err error
	var msg string

	b.logBlockInfo(n, "S")

	// check if the plugin is of type "shell"
	// then include the shell commands all other properties are ignored
	if SHELL == n.Name {
		msg, err = b.invokeShell(n)
	}

	if err != nil {
		return msg, err
	}

	// 1. Check if plugin already available

	return "", nil
}

func (b *builder) logBlockInfo(n *graph.N, when string) {

	if when == "E" || when == "F" {
		log.Println(fmt.Sprintf("%s:~%s:%s:%s:%s~", when, n.ID, n.Name, n.Description, n.Duration))
	} else {
		log.Println(fmt.Sprintf("%s:~%s:%s:%s~", when, n.ID, n.Name, n.Description))
	}
}
