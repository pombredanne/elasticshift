package vcs

import (
	"context"
	"encoding/json"
	"net/http"
)

// AuthorizeResponse ..
type AuthorizeResponse struct {
	Err     error
	URL     string
	Request *http.Request
}

// GetVCSResponse ..
// Used to return any struct of list response
type GetVCSResponse struct {
	Result []VCS
}

func encodeAuthorizeResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {

	resp := r.(AuthorizeResponse)
	if resp.Err != nil {
		return resp.Err
	}

	http.Redirect(w, resp.Request, resp.URL, http.StatusTemporaryRedirect)
	return nil
}

func encodeGetVCSResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {

	resp := r.(GetVCSResponse)
	return json.NewEncoder(w).Encode(resp.Result)
}
