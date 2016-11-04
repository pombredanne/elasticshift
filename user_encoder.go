package esh

import (
	"context"
	"net/http"
	"strings"
	"time"

	chttp "gitlab.com/conspico/esh/core/http"
)

type genericResponse struct {
	Valid bool
	Err   error
}

type signInResponse struct {
	Token string
	Err   error
}

func encodeSignInResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {

	resp := r.(signInResponse)
	if strings.EqualFold(resp.Token, errInvalidEmailOrPassword.Error()) {

		// Setting the error code on encode phase only applicable for
		// sign in request
		w.WriteHeader(http.StatusUnauthorized)
		return resp.Err
	} else if resp.Err != nil {
		return resp.Err
	}

	cookie := &http.Cookie{
		Name:     chttp.AuthTokenCookieName,
		Value:    resp.Token,
		Expires:  time.Now().Add(time.Minute * 15),
		HttpOnly: true,
		Path:     "/",
		//Secure : true, // TODO enable this to ensure the cookie is passed only with https
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)

	return nil
}

func encodeSignOutResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {

	//remove the existing cookie.
	w.Header().Del("Set-Cookie")

	// sets the new dummy cookie.
	cookie := &http.Cookie{
		Name:     chttp.AuthTokenCookieName,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		//Secure : true, // TODO enable this to ensure the cookie is passed only with https
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	return nil
}

func encodeVerifyCodeRequest(ctx context.Context, w http.ResponseWriter, r interface{}) error {

	resp := r.(genericResponse)
	if resp.Err != nil {
		return resp.Err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
