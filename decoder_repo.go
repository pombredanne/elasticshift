// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

import (
	"context"
	"net/http"

	"gitlab.com/conspico/esh/core/auth"
)

// GetRepoRequest ..
type GetRepoRequest struct {
	TeamID string
	VCSID  string
}

func decodeGetRepoRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	team := ctx.Value("token").(auth.Token).Team
	req := GetRepoRequest{TeamID: team}
	params := ctx.Value("params").(map[string]string)
	if params != nil {
		vcsID := params["id"]
		req.VCSID = vcsID
	}
	return req, nil
}
