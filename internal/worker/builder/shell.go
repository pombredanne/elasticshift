/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/keys"
)

func (b *builder) invokeShell(n *N) (string, error) {

	cmds := n.Item()[keys.COMMAND].([]string)

	for _, command := range cmds {

		b.logr.Log(fmt.Sprintf("%s:%s-%s", START, n.Name, n.Description))

		msg, err := b.execShellCmd(n.Name, command, nil, "")
		if err != nil {
			return msg, err
		}
		b.logr.Log(fmt.Sprintf("%s:%s-%s", END, n.Name, n.Description))
	}
	return "", nil
}

func (b *builder) execShellCmd(prefix string, shellCmd string, env []string, dir string) (string, error) {

	cmd := exec.Command("sh", "-c", shellCmd)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	var buf bytes.Buffer
	go io.Copy(b.logr.Writer, stdout)
	go io.Copy(io.MultiWriter(b.logr.Writer, &buf), stderr)

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

	if env != nil {
		cmd.Env = env
	}

	if dir != "" {
		cmd.Dir = dir
	}

	if err := cmd.Start(); err != nil {
		b.log.Errorln(err)
		return buf.String(), err
	}

	if err := cmd.Wait(); err != nil {

		err := fmt.Errorf("Error waiting for the shell command to finish : %v", err)
		b.log.Errorln(err)
		return buf.String(), err
	}

	return "", nil
}
