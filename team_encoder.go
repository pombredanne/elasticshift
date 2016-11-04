package esh

import (
	"context"
	"net/http"
)

type createTeamResponse struct {
	Created bool
	Code    string
	Err     error
}

func encodeCreateTeamResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {

	resp := response.(createTeamResponse)
	if resp.Created {
		w.WriteHeader(http.StatusCreated)
		return nil
	}
	return resp.Err
}
