/*
Copyright 2018 The Elasticshift Authors.
*/
package utils

import (
	"fmt"
	"net"
	"os"
)

func GetIP() string {

	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return ""
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return ""
	}

	// ipv4
	return addrs[0]
}
