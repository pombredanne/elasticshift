/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"bytes"
	"io"
	"os/exec"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/internal/pkg/graph"
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/keys"
)

func (b *builder) invokeShell(n *graph.N) (string, error) {

	cmds := n.Item()[keys.COMMAND].([]string)

	for _, command := range cmds {

		n.Logger.Printf("COMMAND: %s\n", command)

		msg, err := b.execShellCmd(n.Logger, command, nil, "")
		if err != nil {
			n.Logger.Errorf("Failed executing command (%s): %v\n", command, err)
			return msg, err
		}
	}
	return "", nil
}

func (b *builder) execShellCmd(nodelogger *logrus.Entry, shellCmd string, env []string, dir string) (string, error) {

	cmd := exec.Command("sh", "-c", shellCmd)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	var buf bytes.Buffer

	go io.Copy(&CommandWriter{Logger: nodelogger, Type: "I"}, stdout)
	go io.Copy(io.MultiWriter(&CommandWriter{Logger: nodelogger, Type: "E"}, &buf), stderr)

	if env != nil {
		cmd.Env = env
	}

	if dir != "" {
		cmd.Dir = dir
	}

	if err := cmd.Start(); err != nil {
		return buf.String(), err
	}

	if err := cmd.Wait(); err != nil {
		return buf.String(), err
	}

	return "", nil
}
