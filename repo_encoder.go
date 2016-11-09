package esh

import (
	"context"
	"encoding/json"
	"net/http"
)

// GetRepoResponse ..
// Used to return any struct of list response
type GetRepoResponse struct {
	Result []Repo
}

func encodeGetRepoResponse(ctx context.Context, w http.ResponseWriter, r interface{}) error {

	resp := r.(GetRepoResponse)
	return json.NewEncoder(w).Encode(resp.Result)
}
