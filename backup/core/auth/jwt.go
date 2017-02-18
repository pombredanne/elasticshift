package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/palantir/stacktrace"
)

var (
	errUnexpectedSigningMethod = "Unexpected signing method : %v"
	errInvalidIssuer           = errors.New("Token seems to be invalid")
	errTokenExpired            = errors.New("Token is expired")
)

// Token ..
type Token struct {
	UserID   string
	Username string
	Team     string
}

// GenerateToken ..
func GenerateToken(key interface{}, t Token) (string, error) {

	tok := jwt.New(jwt.SigningMethodRS512)
	claims := tok.Claims.(jwt.MapClaims)
	now := time.Now()
	claims["iat"] = now
	claims["exp"] = now.Add(time.Minute * 15).Unix()
	claims["iss"] = "elasticshift.com"
	claims["tok"] = t
	signedString, err := tok.SignedString(key.(*rsa.PrivateKey))
	if err != nil {
		return "", stacktrace.Propagate(err, "Generate token failed.")
	}

	return signedString, err
}

// VefifyToken ..
func VefifyToken(key interface{}, signedToken string) (*jwt.Token, error) {

	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {

		claims := token.Claims.(jwt.MapClaims)
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return false, fmt.Errorf(errUnexpectedSigningMethod, claims["alg"])
		}

		if claims["iss"] != "elasticshift.com" {
			return false, stacktrace.Propagate(errInvalidIssuer, "Verify token failed.")
		}

		// if time.Now().Sub(claims["exp"].(time.Time)) > 0 {
		// 	return false, errTokenExpired
		// }
		return key, nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

// RefreshToken ..
func RefreshToken(key interface{}, token *jwt.Token) (string, error) {

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	signedString, err := token.SignedString(key.(*rsa.PrivateKey))

	if err != nil {
		return "", stacktrace.Propagate(err, "Refresh token failed")
	}

	return signedString, err
}

// GetToken ..
// Parse the JWT token
func GetToken(token *jwt.Token) Token {

	claims := token.Claims.(jwt.MapClaims)
	tok := claims["tok"].(map[string]interface{})

	return Token{
		Team:     tok["Team"].(string),
		UserID:   tok["UserID"].(string),
		Username: tok["Username"].(string),
	}
}
