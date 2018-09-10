/*
Copyright 2018 The Elasticshift Authors.
*/
package worker

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
)

// Generate the RSA keys
// Used to ssh to running containers.
func (w *W) GenerateRSAKeys() error {

	log.Println("Generating RSA keys..")
	r := rand.Reader

	key, err := rsa.GenerateKey(r, DEFAULT_BIT_SIZE)
	if err != nil {
		log.Printf("Failed to generate rsa keys: %v", err)
	}

	sshdir, err := GetSSHDir()
	if err != nil {
		log.Printf("Failed to get ssh dir: %v", err)
	}

	w.privKeyPath = filepath.Join(sshdir, PRIV_KEY_NAME)
	w.pubKeyPath = filepath.Join(sshdir, PUB_KEY_NAME)

	// creates the ssh directory
	utils.Mkdir(sshdir)

	err = w.savePrivateKey(PRIV_KEY_PATH, key)
	if err != nil {
		w.Fatal(fmt.Errorf("Failed to save the private key: %v", err))
	}

	err = w.savePublicKey(PUB_KEY_PATH, key.PublicKey)
	if err != nil {
		w.Fatal(fmt.Errorf("Failed to save the public key: %v", err))
	} else {
		log.Println("Keys generated successfully.")
	}

	return nil
}

// Save the private key in a given filepath
func (w *W) savePrivateKey(filepath string, key *rsa.PrivateKey) error {

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("Failed to save privatekey: %v", err)
	}
	defer f.Close()

	var privatekey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(f, privatekey)
	if err != nil {
		return fmt.Errorf("Failed to PEM encode the private key: %v", err)
	}

	return nil
}

// Save the public key in a given filepath
func (w *W) savePublicKey(filepath string, key rsa.PublicKey) error {

	derEncodedPKIXbytes, err := x509.MarshalPKIXPublicKey(&key)
	if err != nil {
		return fmt.Errorf("Failed to marshall pkix publickey: %v", err)
	}

	var publickey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derEncodedPKIXbytes,
	}

	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("Failed to create %s: %v", filepath, err)
	}
	defer f.Close()

	err = pem.Encode(f, publickey)
	if err != nil {
		return fmt.Errorf("Failed to PEM encode the public key: %v", err)
	}

	return nil
}

// Reads the private key from the given filepath
func (w *W) ReadPrivateKey(filepath string) (string, error) {

	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("Failed to read private key: %v", err)
	}

	return string(b), nil
}
