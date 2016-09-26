package user

import (
	"context"
	"net/http"
	"time"
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
	if resp.Err != nil {
		return resp.Err
	}

	cookie := &http.Cookie{
		Name:     "access-token",
		Value:    resp.Token,
		Expires:  time.Now().Add(time.Minute * 15),
		HttpOnly: true,
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
