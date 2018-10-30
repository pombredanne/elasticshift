/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"github.com/sirupsen/logrus"
)

type CommandWriter struct {
	Logger *logrus.Entry
	Type   string
}

func (cw *CommandWriter) Write(b []byte) (int, error) {

	// data := string(b)
	// if !strings.HasSuffix(data, "\n") {
	// 	data = data + "\n"
	// }

	if cw.Type == "E" {
		cw.Logger.Error(string(b))
	} else {
		cw.Logger.Info(string(b))
	}
	return len(b), nil
}
