package esh

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	chttp "gitlab.com/conspico/esh/core/http"
)

// MakeTeamHandler ..
func MakeTeamHandler(ctx context.Context, s TeamService, r *mux.Router) {

	createTeamHandler := chttp.NewPublicRequestHandler(
		ctx,
		decodeCreateTeamRequest,
		encodeCreateTeamResponse,
		makeCreateTeamEdge(s),
	)

	r.Handle("/api/teams", createTeamHandler).Methods("POST")
}

// MakeUserHandler ..
func MakeUserHandler(ctx context.Context, s UserService, r *mux.Router, signer interface{}, verifier interface{}) {

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

	r.Handle("/api/users/signup", signUpHandler).Methods("POST")
	r.Handle("/api/users/signin", signInHandler).Methods("POST")
	r.Handle("/api/users/signout", signOutHandler).Methods("POST")
	r.Handle("/api/users/verify/{code}", verifyCodeHandler).Methods("POST")
}

// MakeVCSHandler ..
func MakeVCSHandler(ctx context.Context, s VCSService, r *mux.Router, signer interface{}, verifier interface{}) {

	authorizeHandler := chttp.NewPrivateRequestHandler(
		ctx,
		decodeAuthorizeRequest,
		encodeAuthorizeResponse,
		makeAuthorizeEdge(s),
		signer,
		verifier,
	)

	authorizedHandler := chttp.NewPublicRequestHandler(
		ctx,
		decodeAuthorizedRequest,
		encodeAuthorizeResponse,
		makeAuthorizedEdge(s),
	)

	getVCSHandler := chttp.NewPrivateRequestHandler(
		ctx,
		decodeGetVCSRequest,
		encodeGetVCSResponse,
		makeGetVCSEdge(s),
		signer,
		verifier,
	)

	syncVCSHandler := chttp.NewPrivateRequestHandler(
		ctx,
		decodeSyncVCSRequest,
		encodeSyncVCSResponse,
		makeSyncVCSEdge(s),
		signer,
		verifier,
	)
	r.Handle("/api/auth/{provider}", authorizeHandler).Methods("GET")
	r.Handle("/api/auth/{provider}/callback/{id}", authorizedHandler).Methods("GET")
	r.Handle("/api/vcs/sync/{id}", syncVCSHandler).Methods("GET")
	r.Handle("/api/vcs", getVCSHandler).Methods("GET")

}

// MakeRepoHandler ..
func MakeRepoHandler(ctx context.Context, s RepoService, r *mux.Router, signer interface{}, verifier interface{}) {

	getRepoHandler := chttp.NewPrivateRequestHandler(
		ctx,
		decodeGetRepoRequest,
		encodeGetRepoResponse,
		makeGetRepoEdge(s),
		signer,
		verifier,
	)

	r.Handle("/api/repos", getRepoHandler).Methods("GET")
	r.Handle("/api/vcs/{id}/repos", getRepoHandler).Methods("GET")
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
