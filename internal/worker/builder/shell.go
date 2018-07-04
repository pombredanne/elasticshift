/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"fmt"
	"log"
	"os/exec"

	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/keys"
)

func (b *builder) invokeShell(n *N) error {

	cmds := n.Item()[keys.COMMAND].([]string)
	for _, command := range cmds {

		log.Println(fmt.Sprintf("%s:%s-%s", START, n.Name, n.Description))

		err := b.execShellCmd(n.Name, command, nil, "")
		if err != nil {
			return err
		}
		log.Println(fmt.Sprintf("%s:%s-%s", END, n.Name, n.Description))
	}
	return nil
}

func (b *builder) execShellCmd(prefix string, shellCmd string, env []string, dir string) error {

	cmd := exec.Command("sh", "-c", shellCmd)

	cmd.Stdout = newStreamer(prefix)

	// soutpipe, err := cmd.StdoutPipe()
	// if err != nil {
	// 	return err
	// }
	// newStreamer(prefix, soutpipe)

	// serrpipe, err := cmd.StderrPipe()
	// if err != nil {
	// 	return err
	// }
	// newStreamer(prefix, serrpipe)

	// combined out
	// cout, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return err
	// }
	// newStreamer(prefix, bytes.NewReader(cout))
	if env != nil {
		cmd.Env = env
	}

	if dir != "" {
		cmd.Dir = dir
	}

	if err := cmd.Start(); err != nil {

		err := fmt.Errorf("Failed to start the shell command: %v", err)
		log.Println(err)
		return err
	}

	if err := cmd.Wait(); err != nil {

		err := fmt.Errorf("Error waiting for the shell command to finish : %v", err)
		log.Println(err)
		return err
	}

	// fmt.Print(stdout.String())
	return nil
}
