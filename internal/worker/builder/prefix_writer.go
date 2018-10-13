/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

type prefixWriter struct {
	w      io.Writer
	prefix string
}

func newPrefixWriter(prefix string, w io.Writer) prefixWriter {
	return prefixWriter{w, prefix}
}

func (pw prefixWriter) Write(b []byte) (int, error) {

	scanner := bufio.NewScanner(strings.NewReader(string(b)))
	scanner.Split(bufio.ScanLines)

	var buf bytes.Buffer
	for scanner.Scan() {
		buf.WriteString(pw.prefix + scanner.Text() + "\n")
	}
	return pw.w.Write(buf.Bytes())
}
