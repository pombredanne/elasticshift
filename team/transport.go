package team

import (
	"context"

	"github.com/gorilla/mux"
	chttp "gitlab.com/conspico/esh/core/http"
)

// MakeRequestHandler ..
func MakeRequestHandler(ctx context.Context, s Service, r *mux.Router) {

	createTeamHandler := chttp.NewRequestHandler(
		ctx,
		decodeCreateTeamRequest,
		encodeCreateTeamResponse,
		makeCreateTeamEdge(s),
	)

	r.Handle("/teams", createTeamHandler).Methods("POST")
}
