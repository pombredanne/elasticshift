package logwriter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
)

var (
	logDir = "/tmp/shiftlogs"
)

type NodeWriter interface {
	Write(b []byte) (int, error)
}

type nodew struct {
	NodeID string
	file   *os.File
	writer io.Writer
}

func newNodeWriter(nodeid string) (NodeWriter, error) {

	nw := &nodew{NodeID: nodeid}

	exist, err := utils.PathExist(logDir)
	if err != nil {
		return nil, fmt.Errorf("Error Initializing logger: %v", err)
	}

	if !exist {

		err = utils.Mkdir(logDir)
		if err != nil {
			return nil, fmt.Errorf("Error mkdir (%s) : %v", logDir, err)
		}
	}

	f, err := os.OpenFile(filepath.Join(logDir, nodeid), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	nw.file = f
	nw.writer = io.MultiWriter([]io.Writer{f, os.Stdout}...)

	return nw, nil
}

func (w *nodew) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}
