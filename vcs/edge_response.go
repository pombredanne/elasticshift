package vcs

import (
	"context"
	"net/http"
)

// AuthorizeResponse ..
type AuthorizeResponse struct {
	Err     error
	URL     string
	Request *http.Request
}

func encodeAuthorizeResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {

	resp := r.(AuthorizeResponse)
	if resp.Err != nil {
		return resp.Err
	}

	if resp.Err != nil {
		http.Redirect(w, resp.Request, resp.URL, http.StatusTemporaryRedirect)
	}
	return nil
}
