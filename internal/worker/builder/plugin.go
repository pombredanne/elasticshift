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

	b.logBlockInfo(n, graph.START)

	// check if the plugin is of type "shell"
	// then include the shell commands all other properties are ignored
	if SHELL == n.Name {
		msg, err = b.invokeShell(n)
	}

	if err != nil {
		log.Println(fmt.Sprintf("%s:~%s-%s: %v~", graph.ERROR, n.Name, n.Description, err))
		return msg, err
	}

	// 1. Check if plugin already available

	b.logBlockInfo(n, graph.END)
	return "", nil
}

func (b *builder) logBlockInfo(n *graph.N, when string) {

	if n.Description == "" {
		log.Println(fmt.Sprintf("%s:~%s~", when, n.Name))
	} else {
		log.Println(fmt.Sprintf("%s:~%s-%s~", when, n.Name, n.Description))
	}
}
