// Package util
// Author Ghazni Nattarshah
// Date: 1/3/17
package util

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

// LoadKey ..
func LoadKey(path string) (interface{}, error) {

	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	keyBlock, _ := pem.Decode(keyBytes)

	switch keyBlock.Type {
	case "PUBLIC KEY":
		return x509.ParsePKIXPublicKey(keyBlock.Bytes)
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	default:
		return nil, fmt.Errorf("unsupported key type %q", keyBlock.Type)
	}
}
