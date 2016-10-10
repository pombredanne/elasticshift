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
		encodeListVCSResponse,
		makeAuthorizedEdge(s),
	)

	listVCSHandler := chttp.NewPrivateRequestHandler(
		ctx,
		decodeListVCSRequest,
		encodeListVCSResponse,
		makeAuthorizeEdge(s),
		signer,
		verifier,
	)

	r.Handle("/api/auth/{provider}", authorizeHandler)
	r.Handle("/api/auth/{provider}/callback/{team}", authorizedHandler)
	r.Handle("/api/vcs/", listVCSHandler)
}
