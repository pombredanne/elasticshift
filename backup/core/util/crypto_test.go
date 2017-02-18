package util_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"gitlab.com/conspico/esh/core/util"
)

var key string = "53b7ef54af124d0d435bba53ce602498"

type data struct {
	Team   string
	Email  string
	Expire time.Time
}

func TestEncrypt(t *testing.T) {

	text := "conspico:ghazni.nattarshah@gmail.com;"
	enc, _ := util.Encrypt(key, []byte(text))
	t.Log("Encrypted = ", enc)

	dec, _ := util.Decrypt(key, enc)
	t.Log("Decrypted = ", dec)
	t.Log("Decrypted = ", fmt.Sprintf("%x", enc))
}

func TestEncryptStruct(t *testing.T) {

	origData := &data{
		Team:   "conapico",
		Email:  "ghazni.nattarshah@conspico.com",
		Expire: time.Now().AddDate(0, 0, 1),
	}

	cipherText, err := util.EncryptStruct(key, origData)
	if err != nil {
		t.Log("Encryption error  = ", err)
		t.FailNow()
	} else {
		t.Log("Encrypted = ", cipherText)
		t.Log("Length of encrypted = ", len(cipherText))
	}
}

func TestDecrypt(t *testing.T) {

	cipherText := "w9m9Y6N2Zl3gZLOq+WcyMKgabaL6QjCJQxa/G2iPeg6VavCacRsfTrRF5DFETcQdUfNR/UkdgRlW2D3ijOSTGsxUH997dERWsSzh2EQB9rbK5R+OHQTB/nBR40MnPA/8CzhtsrO/jmSvLfQAAAAAAAAAAAAAAAAAAAAA"
	var decryptedData data
	util.DecryptStruct(key, cipherText, decryptedData)
	t.Log("Decrypted = ", decryptedData)
}

func TestSha512Hash(t *testing.T) {

	text := "conspico:ghazni.nattarshah@conspico.com"
	hashed := util.Sha512Hash(text)
	t.Log("Hashed = ", hashed)

	isValid := util.CompareSha512Hash(text, hashed)
	t.Log("Valid = ", isValid)
}

func TestXOREncryptDecrypt(t *testing.T) {

	id, _ := util.NewUUID()
	//text := id + ":conspico:ghazni.nattarshah@conspico.com"

	var v bytes.Buffer
	v.WriteString(id)
	v.WriteString(";")
	v.WriteString("ghazni.nattarshah@conspico.com;")
	expireAt, _ := time.Now().AddDate(0, 0, 7).MarshalText()
	v.Write(expireAt)

	enc := util.XOREncrypt(key, v.String())
	t.Log("XOREncrypted = ", enc)

	dec, err := util.XORDecrypt(key, enc)
	if err != nil {
		t.Log("XORDecryption error = ", err)
	}
	t.Log("XORDecrypted = ", dec)
}
