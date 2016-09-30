package vcs

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

// AuthorizeRequest ..
type AuthorizeRequest struct {
	Provider string
	Domain   string
	Request  *http.Request
}

func decodeAuthorizeRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	subdomain := ctx.Value("subdomain").(string)

	params := mux.Vars(r)
	prov := params["provider"]

	return AuthorizeRequest{Domain: subdomain, Provider: prov, Request: r}, nil
}
