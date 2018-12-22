/*
Copyright 2018 The Elasticshift Authors.
*/
package main

import (
	"os"

	"log"

	"github.com/elasticshift/elasticshift/internal/worker"
)

func main() {

	if err := worker.Run(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
