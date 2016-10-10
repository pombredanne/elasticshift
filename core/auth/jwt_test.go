package auth_test

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"testing"

	"gitlab.com/conspico/esh/core/auth"
)

func TestJWT(t *testing.T) {

	// load keys
	signerBytes, err := ioutil.ReadFile("../../keys/esh.rsa")
	signerBlock, _ := pem.Decode(signerBytes)
	signer, err := x509.ParsePKCS1PrivateKey(signerBlock.Bytes)
	if err != nil {
		t.Log("Failed to load signer key", err)
	}
	verifierBytes, err := ioutil.ReadFile("../../keys/esh.rsa.pub")
	verifierBlock, _ := pem.Decode(verifierBytes)
	verifier, err := x509.ParsePKIXPublicKey(verifierBlock.Bytes)

	if err != nil {
		t.Log("Failed to load verfier key", err)
	}

	tok := auth.Token{UserID: "ghazni.nattarshah@conspico.com", TeamID: "conspico.com"}
	signedtoken, err := auth.GenerateToken(signer, tok)
	if err != nil {
		t.Log("Failed to generate auth token : ", err)
	}

	verifiedToken, err := auth.VefifyToken(verifier, signedtoken)
	if err != nil {
		t.Log("Failed to verify token", err)
	}
	if verifiedToken != nil {
		t.Log("Token valid = ", verifiedToken.Valid)
	}

	origTok := auth.GetToken(verifiedToken)
	t.Log(origTok.TeamID)
}
