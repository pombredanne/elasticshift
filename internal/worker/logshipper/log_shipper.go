/*
Copyright 2018 The Elasticshift Authors.
*/
package logshipper

import (
	"bufio"
	"os"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/internal/pkg/storage"
)

type LogShipper interface {
	Ship(f *os.File)
}

type logshipper struct {
	queue   chan *os.File
	logger  *logrus.Entry
	storage *storage.ShiftStorage
}

// New ..
// Creates a new log shipper
func New(l *logrus.Entry, s *storage.ShiftStorage) (LogShipper, error) {

	ls := &logshipper{
		queue:   make(chan *os.File),
		storage: s,
		logger:  l,
	}

	go ls.start()

	return ls, nil
}

func (s *logshipper) Ship(f *os.File) {
	s.queue <- f
}

func (s *logshipper) start() {

	for f := range s.queue {

		err := s.storage.PutLog(f.Name(), bufio.NewReader(f))
		if err != nil {
			s.logger.Errorf("Failed to store build log %v \n ", err)
		}
	}
}
