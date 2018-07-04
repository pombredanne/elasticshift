/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

var (
	SHELL = "elasticshift/shell"
)

func (b *builder) invokePlugin(n *N) error {

	if START == n.Name || END == n.Name {
		return nil
	}

	var err error
	// check if the plugin is o type "elasticshift/shell"
	// then include the shell commands all other properties are ignored
	if SHELL == n.Name {
		err = b.invokeShell(n)
	}

	if err != nil {
		return err
	}

	// 1. Check if plugin already available

	return nil
}
