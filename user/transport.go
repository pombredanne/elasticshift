package user

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	chttp "gitlab.com/conspico/esh/core/http"
)

// MakeRequestHandler ..
func MakeRequestHandler(ctx context.Context, s Service, r *mux.Router, signer interface{}, verifier interface{}) {

	signUpHandler := chttp.NewPublicRequestHandler(
		ctx,
		decodeSignUpRequest,
		encodeSignInResponse,
		makeSignupEdge(s),
	)

	signInHandler := chttp.NewPublicRequestHandler(
		ctx,
		decodeSignInRequest,
		encodeSignInResponse,
		makeSignInEdge(s),
	)

	signOutHandler := chttp.NewPrivateRequestHandler(
		ctx,
		decodeSignOutRequest,
		encodeSignOutResponse,
		makeSignOutEdge(s),
		signer,
		verifier,
	)

	verifyCodeHandler := chttp.NewPublicRequestHandler(
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

	r.Handle("/api/users/signup", accessControl(signUpHandler)).Methods("POST")
	r.Handle("/api/users/signin", accessControl(signInHandler)).Methods("POST")
	r.Handle("/api/users/signout", accessControl(signOutHandler)).Methods("POST")
	r.Handle("/api/users/verify/{code}", verifyCodeHandler).Methods("POST")
	//r.Handle("/users/verifyAndSignIn", verifyAndSignInHandler).Methods("POST")
	//r.Handle("/users/signin", signinHandler).Methods("POST")
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
