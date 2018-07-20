/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"bytes"
	"io"
)

type streamer struct {
	w      io.Writer
	buf    bytes.Buffer
	stderr bool
}

func newStreamer(w io.Writer, stderr bool) streamer {

	s := streamer{}
	s.w = w
	if stderr {
		s.buf = bytes.Buffer{}
	}
	return s
}

func (s streamer) Write(b []byte) (int, error) {

	if s.stderr {
		s.buf.Write(b)
	}
	return s.w.Write(b)
}

func (s streamer) Error() string {
	return s.buf.String()
}
