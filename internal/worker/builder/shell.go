/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os/exec"

	"gitlab.com/conspico/elasticshift/internal/pkg/graph"
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/keys"
)

func (b *builder) invokeShell(n *graph.N) (string, error) {

	cmds := n.Item()[keys.COMMAND].([]string)

	var prefix, encmd string
	for _, command := range cmds {

		encmd = base64.StdEncoding.EncodeToString([]byte(command))
		log.Println(fmt.Sprintf("S:EXEC:~%s:%s~", n.ID, encmd))

		prefix = ""
		if n.Parallel {
			prefix = n.ID + "[!@#$]"
		}

		msg, err := b.execShellCmd(prefix, command, nil, "")
		if err != nil {
			log.Println(fmt.Sprintf("F:EXEC:~%s:%s~", n.ID, encmd))
			return msg, err
		}
		log.Println(fmt.Sprintf("E:EXEC:~%s:%s~", n.ID, encmd))
	}
	return "", nil
}

func (b *builder) execShellCmd(prefix string, shellCmd string, env []string, dir string) (string, error) {

	cmd := exec.Command("sh", "-c", shellCmd)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	var buf bytes.Buffer
	var w io.Writer

	if prefix != "" {
		w = newPrefixWriter(prefix, b.writer)
	} else {
		w = b.writer
	}

	go io.Copy(w, stdout)
	go io.Copy(io.MultiWriter(w, &buf), stderr)

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
