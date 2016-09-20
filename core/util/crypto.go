package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
)

// Encrypt - plain text to encoded string
func Encrypt(key string, text []byte) (string, error) {

	// creates a cipher block with key bytes
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// creates a random iv and have it first part of cipher text
	ciperText := make([]byte, aes.BlockSize+len(text))
	iv := ciperText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {

	}

	// create a encrypter with iv
	cfb := cipher.NewCFBEncrypter(block, iv)

	cfb.XORKeyStream(ciperText, text)

	return base64.StdEncoding.EncodeToString(ciperText), nil
}

// Decrypt - cipher to text
func Decrypt(key, cipherText string) ([]byte, error) {

	// decodes the cipher text
	cipherBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}

	// creates a cipher block with key bytes
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// Extract iv from cipher bytes
	iv := cipherBytes[:aes.BlockSize]

	// create a decrypter with iv
	cbs := cipher.NewCFBDecrypter(block, iv)

	text := cipherBytes[aes.BlockSize:]
	cbs.XORKeyStream(text, text)

	// convert bytes to plaintext
	return text, nil
}

// EncryptStruct ..
func EncryptStruct(key string, value interface{}) (string, error) {

	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return Encrypt(key, b)
}

// DecryptStruct ..
func DecryptStruct(key string, cipherText string, value interface{}) error {

	b, err := Decrypt(key, cipherText)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &value)
}
