package team

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	chttp "gitlab.com/conspico/esh/core/http"
)

// MakeRequestHandler ..
func MakeRequestHandler(ctx context.Context, s Service) http.Handler {

	createTeamHandler := chttp.NewRequestHandler(
		ctx,
		decodeCreateTeamRequest,
		encodeCreateTeamResponse,
		makeCreateTeamEdge(s),
	)

	r := mux.NewRouter()
	r.Handle("/teams/v1", createTeamHandler).Methods("POST")

	return r
}
