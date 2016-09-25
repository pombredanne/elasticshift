package user

import (
	"context"
	"net/http"
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
