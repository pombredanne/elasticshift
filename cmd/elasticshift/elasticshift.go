/*
Copyright 2017 The Elasticshift Authors.
*/
package main

import (
	"os"

	"log"

	"gitlab.com/conspico/elasticshift/internal/shiftserver"
)

func main() {

	if err := shiftserver.Run(); err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
