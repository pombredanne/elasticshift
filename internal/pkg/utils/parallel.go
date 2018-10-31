/*
Copyright 2018 The Elasticshift Authors.
*/
package utils

import "runtime"

func NumOfCPU() int {

	var parallel int
	nCpu := runtime.NumCPU()
	if nCpu < 2 {
		parallel = 1
	} else {
		parallel = nCpu - 1
	}
	return parallel
}
