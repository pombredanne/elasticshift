package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/palantir/stacktrace"
)

// Encrypt - plain text to encoded string
func Encrypt(key string, text []byte) (string, error) {

	// creates a cipher block with key bytes
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", stacktrace.Propagate(err, "Can't encrypt")
	}

	// creates a random iv and have it first part of cipher text
	cipherText := make([]byte, aes.BlockSize+len(text))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {

	}

	// create a encrypter with iv
	cfb := cipher.NewCFBEncrypter(block, iv)

	cfb.XORKeyStream(cipherText, text)

	return base64.StdEncoding.EncodeToString(cipherText), nil
	//return hex.EncodeToString(cipherText[:]), nil
}

// Decrypt - cipher to text
func Decrypt(key, cipherText string) ([]byte, error) {

	// decodes the cipher text
	cipherBytes, err := base64.StdEncoding.DecodeString(cipherText)
	//cipherBytes, err := hex.DecodeString(cipherText)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Can't decrypt")
	}

	// creates a cipher block with key bytes
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, stacktrace.Propagate(err, "Can't decrypt")
	}

	//cipherBytes := []byte(cipherText)
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
		return "", stacktrace.Propagate(err, "Can't encrypt")
	}
	return Encrypt(key, b)
}

// DecryptStruct ..
func DecryptStruct(key string, cipherText string, value interface{}) error {

	b, err := Decrypt(key, cipherText)
	if err != nil {
		return stacktrace.Propagate(err, "Can't decrypt")
	}

	return json.Unmarshal(b, &value)
}

// Sha512Hash conversion
func Sha512Hash(text string) string {

	hashed := sha512.Sum512([]byte(text))
	//return base64.StdEncoding.EncodeToString(hashed[:])
	//return string(hashed[:])
	return hex.EncodeToString(hashed[:])
}

// CompareSha512Hash ..
func CompareSha512Hash(plain, hashed string) bool {
	return hashed == Sha512Hash(plain)
}

// XOREncrypt ..
func XOREncrypt(key, input string) string {

	enc := xorEncrypeDecrypt(key, input)
	return hex.EncodeToString([]byte(enc))
	//return enc
}

// XORDecrypt ..
func XORDecrypt(key, input string) (string, error) {

	decoded, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}
	return xorEncrypeDecrypt(key, string(decoded[:])), nil
}

func xorEncrypeDecrypt(key, input string) string {

	var output string
	for i := 0; i < len(input); i++ {
		output += string(input[i] ^ key[i%len(key)])
	}
	return output
}
