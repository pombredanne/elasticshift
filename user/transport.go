package user

import (
	"context"

	"github.com/gorilla/mux"
	chttp "gitlab.com/conspico/esh/core/http"
)

// MakeRequestHandler ..
func MakeRequestHandler(ctx context.Context, s Service, r *mux.Router) {

	createUserHandler := chttp.NewRequestHandler(
		ctx,
		decodeCreateUserRequest,
		encodeCreateUserResponse,
		makeCreateUserEdge(s),
	)

	verifyCodeHandler := chttp.NewRequestHandler(
		ctx,
		decodeVerifyCodeRequest,
		encodeVerifyCodeRequest,
		makeVerifyCodeEdge(s),
	)

	// verifyAndSignInHandler := chttp.NewRequestHandler(
	// 	ctx,
	// 	decodeCreateUserRequest,
	// 	encodeCreateUserResponse,
	// 	makeCreateUserEdge(s),
	// )

	// signinHandler := chttp.NewRequestHandler(
	// 	ctx,
	// 	decodeCreateUserRequest,
	// 	encodeCreateUserResponse,
	// 	makeCreateUserEdge(s),
	// )

	r.Handle("/users", createUserHandler).Methods("POST")
	r.Handle("/users/verify/{code}", verifyCodeHandler).Methods("POST")
	//r.Handle("/users/verifyAndSignIn", verifyAndSignInHandler).Methods("POST")
	//r.Handle("/users/signin", signinHandler).Methods("POST")

}
