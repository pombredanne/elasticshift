/*
Copyright 2018 The Elasticshift Authors.
*/
package logshipper

import (
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/internal/pkg/storage"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
)

type LogShipper interface {
	Ship(nodeid, filepath string)
	WaitUntilLogShipperCompletes()
}

type logshipper struct {
	queue   chan *info
	logger  *logrus.Entry
	storage *storage.ShiftStorage
	wg      sync.WaitGroup
}

type info struct {
	nodeid   string
	filepath string
}

// New ..
// Creates a new log shipper
func New(l *logrus.Entry, s *storage.ShiftStorage) (LogShipper, error) {

	ls := &logshipper{
		queue:   make(chan *info),
		storage: s,
		logger:  l,
	}

	go ls.start()

	return ls, nil
}

func (s *logshipper) Ship(nodeid, filepath string) {
	s.wg.Add(1)
	s.queue <- &info{nodeid, filepath}
}

func (s *logshipper) start() {

	parallelCh := make(chan int, utils.NumOfCPU())
	for i := range s.queue {

		go func(i *info) {

			parallelCh <- 1
			defer s.wg.Done()

			err := s.storage.PutLog(i.nodeid, i.filepath)
			if err != nil {
				s.logger.Errorf("Failed to store build log %v \n ", err)
			}

			<-parallelCh
		}(i)
	}
}

func (s *logshipper) WaitUntilLogShipperCompletes() {
	s.wg.Wait()
}
