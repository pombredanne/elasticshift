package vcs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// AuthorizeRequest ..
type AuthorizeRequest struct {
	Provider string
	Domain   string
	Request  *http.Request
	Code     string
}

// GetVCSRequest ..
type GetVCSRequest struct {
	Domain string
}

func decodeAuthorizeRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	subdomain := ctx.Value("subdomain").(string)

	params := mux.Vars(r)
	prov := params["provider"]

	return AuthorizeRequest{Domain: subdomain, Provider: prov, Request: r}, nil
}

func decodeAuthorizedRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	params := mux.Vars(r)
	prov := params["provider"]
	subdomain := params["team"]
	code := r.FormValue("code")

	fmt.Println("Subdomain from callback =", subdomain)
	return AuthorizeRequest{Domain: subdomain, Provider: prov, Request: r, Code: code}, nil
}

func decodeGetVCSRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	subdomain := ctx.Value("subdomain").(string)
	return subdomain, nil
}
