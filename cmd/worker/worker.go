/*
Copyright 2018 The Elasticshift Authors.
*/
package main

import (
	"os"

	"log"

	"gitlab.com/conspico/elasticshift/internal/worker"
)

func main() {

	if err := worker.Run(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
