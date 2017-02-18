// Package esh ...
// Author: Ghazni Nattarshah
// Date: Sep 24, 2016
package esh

import (
	"context"
	"net/http"

	"gitlab.com/conspico/esh/core/auth"
)

// AuthorizeRequest ..
type AuthorizeRequest struct {
	Provider string
	Team     string
	ID       string
	Request  *http.Request
	Code     string
}

// SyncVCSRequest ..
type SyncVCSRequest struct {
	Team       string
	Username   string
	ProviderID string
}

// GetVCSRequest ..
type GetVCSRequest struct {
	Domain string
}

func decodeAuthorizeRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	team := ctx.Value("token").(auth.Token).Team
	params := ctx.Value("params").(map[string]string)
	prov := params["provider"]

	return AuthorizeRequest{Team: team, Provider: prov, Request: r}, nil
}

func decodeAuthorizedRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	params := ctx.Value("params").(map[string]string)
	prov := params["provider"]
	id := r.FormValue("id")
	code := r.FormValue("code")

	return AuthorizeRequest{ID: id, Provider: prov, Request: r, Code: code}, nil
}

func decodeGetVCSRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return ctx.Value("token").(auth.Token).Team, nil
}

func decodeSyncVCSRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	params := ctx.Value("params").(map[string]string)
	providerID := params["id"]
	token := ctx.Value("token").(auth.Token)

	return SyncVCSRequest{Team: token.Team, Username: token.Username, ProviderID: providerID}, nil
}
