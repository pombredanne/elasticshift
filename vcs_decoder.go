package esh

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/conspico/esh/core/auth"

	"github.com/gorilla/mux"
)

// AuthorizeRequest ..
type AuthorizeRequest struct {
	Provider string
	TeamID   string
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

	params := mux.Vars(r)
	prov := params["provider"]

	return AuthorizeRequest{TeamID: teamID, Provider: prov, Request: r}, nil
}

func decodeAuthorizedRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	params := mux.Vars(r)
	prov := params["provider"]
	teamID := params["team"]
	code := r.FormValue("code")

	fmt.Println("Subdomain from callback =", teamID)
	return AuthorizeRequest{TeamID: teamID, Provider: prov, Request: r, Code: code}, nil
}

func decodeGetVCSRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	teamID := ctx.Value("token").(auth.Token).TeamID
	fmt.Println("TeamID = ", teamID)
	return teamID, nil
}

func decodeSyncVCSRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	params := mux.Vars(r)
	providerID := params["id"]
	token := ctx.Value("token").(auth.Token)

	return SyncVCSRequest{TeamID: token.TeamID, Username: token.Username, ProviderID: providerID}, nil
}
