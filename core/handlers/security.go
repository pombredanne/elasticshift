package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"gitlab.com/conspico/esh/core/auth"
)

var (
	unAuthorized    = "Unauthorized"
	errUnauthorized = errors.New(unAuthorized)

	// AuthTokenCookieName ..
	AuthTokenCookieName = "__at"
)

type security struct {
	ctx      context.Context
	h        http.Handler
	signer   interface{}
	verifier interface{}
}

func (sh *security) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	cookie, err := r.Cookie(AuthTokenCookieName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, unAuthorized, http.StatusUnauthorized)
		return
	}

	if "" == cookie.Value {
		fmt.Println(err)
		http.Error(w, unAuthorized, http.StatusUnauthorized)
		return
	}

	token, err := auth.VefifyToken(sh.verifier, cookie.Value)
	if err != nil || !token.Valid {
		fmt.Println(err)
		http.Error(w, unAuthorized, http.StatusUnauthorized)
		return
	}

	c := context.WithValue(ctx, "token", auth.GetToken(token))

	// Refresh the token
	refreshtoken(token, sh.signer, w)

	sh.h.ServeHTTP(w, r.WithContext(c))
}

func refreshtoken(token *jwt.Token, signer interface{}, w http.ResponseWriter) {

	signedTok, err := auth.RefreshToken(signer, token)
	if err != nil {
		fmt.Println("Failed to refresh the token.", err)
	}

	cookie := &http.Cookie{
		Name:     AuthTokenCookieName,
		Value:    signedTok,
		Expires:  time.Now().Add(time.Minute * 15),
		HttpOnly: true,
		Path:     "/",
		//Secure : true, // TODO enable this to ensure the cookie is passed only with https
	}
	http.SetCookie(w, cookie)
}

// SecurityHandler ..
func SecurityHandler(ctx context.Context, signer interface{}, verifier interface{}) func(http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {

		return &security{
			ctx:      ctx,
			signer:   signer,
			verifier: verifier,
			h:        h,
		}
	}
}
