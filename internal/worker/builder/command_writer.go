/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import "github.com/Sirupsen/logrus"

type CommandWriter struct {
	Logger *logrus.Entry
	Type   string
}

func (cw *CommandWriter) Write(b []byte) (int, error) {
	if cw.Type == "E" {
		cw.Logger.Error(string(b))
	} else {
		cw.Logger.Info(string(b))
	}
	return len(b), nil
}
