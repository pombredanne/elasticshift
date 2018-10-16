/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
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

	// check if the plugin is of type "shell"
	// then include the shell commands all other properties are ignored
	if SHELL == n.Name {
		msg, err = b.invokeShell(n)
	} else if graph.RESTORE_CACHE == n.Name {
		err = b.restoreCache(n.Logger)
	} else if graph.SAVE_CACHE == n.Name {
		err = b.saveCache(n.Logger)
	}

	if err != nil {
		return msg, err
	}

	// 1. Check if plugin already available

	return "", nil
}
