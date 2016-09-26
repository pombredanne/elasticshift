package team

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	chttp "gitlab.com/conspico/esh/core/http"
)

// MakeRequestHandler ..
func MakeRequestHandler(ctx context.Context, s Service, r *mux.Router) {

	createTeamHandler := chttp.NewPublicRequestHandler(
		ctx,
		decodeCreateTeamRequest,
		encodeCreateTeamResponse,
		makeCreateTeamEdge(s),
	)

	r.Handle("/api/teams", accessControl(createTeamHandler)).Methods("POST")
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
