/*
Copyright 2019 The Elasticshift Authors.
*/
package testhelper

import "github.com/elasticshift/elasticshift/internal/pkg/logger"

// GetLoggr ...
// Gets the test logger
func GetLoggr() (logger.Loggr, error) {
	return logger.New("info", "text")
}
