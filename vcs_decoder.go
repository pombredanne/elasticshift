package esh

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/conspico/esh/core/auth"
)

// AuthorizeRequest ..
type AuthorizeRequest struct {
	Provider string
	TeamID   string
	ID       string
	Request  *http.Request
	Code     string
}

// SyncVCSRequest ..
type SyncVCSRequest struct {
	TeamID     string
	Username   string
	ProviderID string
}

// GetVCSRequest ..
type GetVCSRequest struct {
	Domain string
}

func decodeAuthorizeRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	teamID := ctx.Value("token").(auth.Token).TeamID
	params := ctx.Value("params").(map[string]string)
	prov := params["provider"]

	return AuthorizeRequest{TeamID: teamID, Provider: prov, Request: r}, nil
}

func decodeAuthorizedRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	params := ctx.Value("params").(map[string]string)
	prov := params["provider"]
	id := params["id"]
	code := r.FormValue("code")

	fmt.Println("Subdomain from callback =", id)
	return AuthorizeRequest{ID: id, Provider: prov, Request: r, Code: code}, nil
}

func decodeGetVCSRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	teamID := ctx.Value("token").(auth.Token).TeamID
	return teamID, nil
}

func decodeSyncVCSRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	params := ctx.Value("params").(map[string]string)
	providerID := params["id"]
	token := ctx.Value("token").(auth.Token)

	return SyncVCSRequest{TeamID: token.TeamID, Username: token.Username, ProviderID: providerID}, nil
}
