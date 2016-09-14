package user

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	chttp "gitlab.com/conspico/esh/core/http"
)

// MakeRequestHandler ..
func MakeRequestHandler(ctx context.Context, s Service) http.Handler {

	createuserHandler := chttp.NewRequestHandler(
		ctx,
		decodeCreateUserRequest,
		encodeCreateUserResponse,
		makeCreateUserEdge(s),
	)

	r := mux.NewRouter()
	r.Handle("/users/v1", createuserHandler).Methods("POST")

	return r
}
