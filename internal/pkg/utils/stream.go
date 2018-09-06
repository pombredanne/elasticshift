/*
Copyright 2018 The Elasticshift Authors.
*/
package utils

import (
	"net/http"
)

type Stream struct {
	buffer int                 // Buffer size
	w      http.ResponseWriter // Underlying writer to send data to
	f      http.Flusher
}

// StreamWriter ..
// Write and flush the data to http response writer
func StreamWriter(w http.ResponseWriter) *Stream {
	s := new(Stream)
	s.w = w
	s.buffer = 1024
	if f, ok := w.(http.Flusher); ok {
		s.f = f
	}
	return s
}

// Write the steam to http response writer and flush it
func (s *Stream) Write(p []byte) (n int, err error) {

	n, err = s.w.Write(p)
	if err != nil {
		return n, err
	}

	if s != nil && s.f != nil {
		s.f.Flush()
	}

	return n, nil
}
