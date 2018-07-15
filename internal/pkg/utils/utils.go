/*
Copyright 2018 The Elasticshift Authors.
*/
package utils

import (
	"log"
	"os"
)

func Mkdir(path string) error {
	exist, err := PathExist(path)
	if err != nil {
		return err
	}
	if !exist {
		os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func GetWD() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
