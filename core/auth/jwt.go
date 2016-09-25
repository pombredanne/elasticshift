package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	errUnexpectedSigningMethod = "Unexpected signing method : %v"
	errInvalidIssuer           = errors.New("Invalid token")
)

// Token ..
type Token struct {
	TeamID string
	Email  string
}

// GenerateToken ..
func GenerateToken(key []byte, t Token) (string, error) {

	tok := jwt.New(jwt.SigningMethodHS512)
	claims := tok.Claims.(jwt.MapClaims)
	now := time.Now()
	claims["iat"] = now
	claims["exp"] = now.Add(time.Minute * 15).Unix()
	claims["iss"] = "elasticshift.com"
	claims["tok"] = t
	signedString, err := tok.SignedString(key)
	if err != nil {
		return "", err
	}
	return signedString, err
}

// VefifyToken ..
func VefifyToken(key []byte, signedToken string) (bool, error) {

	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {

		claims := token.Claims.(jwt.MapClaims)
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return false, fmt.Errorf(errUnexpectedSigningMethod, claims["alg"])
		}

		if claims["iss"] != "elasticshift.com" {
			return false, errInvalidIssuer
		}
		return key, nil
	})

	if err != nil {
		return false, err
	}
	return token.Valid, nil
}
