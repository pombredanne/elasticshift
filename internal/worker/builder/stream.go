/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"log"
)

type streamer struct {
	prefix string
}

func newStreamer(prefix string) streamer {

	s := streamer{}
	s.prefix = prefix
	return s
}

func (s streamer) Write(b []byte) (int, error) {

	l := len(b)
	log.Printf("%s: %s", s.prefix, b)
	return l, nil
}
