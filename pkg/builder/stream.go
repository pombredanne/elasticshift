/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

type streamer struct {
	name    string
	scanner *bufio.Scanner
}

func newStreamer(name string, reader io.Reader) streamer {

	s := streamer{}
	s.name = name

	s.scanner = bufio.NewScanner(reader)

	go s.stream()

	return s
}

func (s *streamer) stream() {

	for s.scanner.Scan() {
		log.Printf("%s: %s\n", s.name, s.scanner.Text())
	}

	err := s.scanner.Err()
	if err != nil {
		log.Println(fmt.Errorf("Error when streaming logs: %v", err))
	}
}
