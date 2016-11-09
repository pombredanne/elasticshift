package esh

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/conspico/esh/core/auth"
)

// GetRepoRequest ..
type GetRepoRequest struct {
	TeamID string
	VCSID  string
}

func decodeGetRepoRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	teamID := ctx.Value("token").(auth.Token).TeamID

	req := GetRepoRequest{TeamID: teamID}
	params := mux.Vars(r)
	if params != nil {
		vcsID := params["id"]
		req.VCSID = vcsID
	}
	return req, nil
}
