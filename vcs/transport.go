package vcs

import (
	"context"

	"github.com/gorilla/mux"
	chttp "gitlab.com/conspico/esh/core/http"
)

// MakeRequestHandler ..
func MakeRequestHandler(ctx context.Context, s Service, r *mux.Router, signer interface{}, verifier interface{}) {

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

	listVCSHandler := chttp.NewPrivateRequestHandler(
		ctx,
		decodeGetVCSRequest,
		encodeGetVCSResponse,
		makeGetVCSEdge(s),
		signer,
		verifier,
	)

	r.Handle("/api/auth/{provider}", authorizeHandler).Methods("GET")
	r.Handle("/api/auth/{provider}/callback/{team}", authorizedHandler).Methods("GET")
	r.Handle("/api/vcs", listVCSHandler).Methods("GET")
}
