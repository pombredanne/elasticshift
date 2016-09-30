package vcs

import (
	"context"

	"github.com/gorilla/mux"
	chttp "gitlab.com/conspico/esh/core/http"
)

// MakeRequestHandler ..
func MakeRequestHandler(ctx context.Context, s Service, r *mux.Router, verifier []byte) {

	authorizeHandler := chttp.NewPrivateRequestHandler(
		ctx,
		decodeAuthorizeRequest,
		encodeAuthorizeResponse,
		makeAuthorizeEdge(s),
		verifier,
	)

	authorizedHandler := chttp.NewPublicRequestHandler(
		ctx,
		decodeAuthorizeRequest,
		encodeAuthorizeResponse,
		makeAuthorizeEdge(s),
	)

	r.Handle("/auth/{provider}", authorizeHandler)
	r.Handle("/auth/{provider}/callback", authorizedHandler)
}
