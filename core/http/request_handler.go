package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"errors"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/conspico/esh/core/auth"
	"gitlab.com/conspico/esh/core/edge"
)

// Enum that represents the PHASE of the request
const (
	AUTH    = 0
	DECODE  = 1
	PROCESS = 2
	ENCODE  = 3
)

var (
	errUnauthorized = errors.New("Unauthorized")

	// AuthTokenCookieName ..
	AuthTokenCookieName = "__at"
)

// RequestHandler ..
// Any request reaches to ESH server lands here
// and a life cycle will be performed such as decode, process, encode  a request
type RequestHandler struct {
	ctx       context.Context
	decode    RequestDecoderFunc
	encode    ResponseEncoderFunc
	process   edge.Edge
	protected bool
	signer    interface{}
	verifier  interface{}
}

// NewPrivateRequestHandler creates a reqeust handler for given edge
// It's a private handler, only authorized request are allowed
func NewPrivateRequestHandler(
	ctx context.Context,
	decoder RequestDecoderFunc,
	encoder ResponseEncoderFunc,
	exec edge.Edge,
	signer interface{},
	verifier interface{}) *RequestHandler {

	rh := &RequestHandler{
		ctx,
		decoder,
		encoder,
		exec,
		true,
		signer,
		verifier,
	}
	return rh
}

// NewPublicRequestHandler creates a reqeust handler for given edge
// It's a private handler, only authorized request are allowed
func NewPublicRequestHandler(
	ctx context.Context,
	decoder RequestDecoderFunc,
	encoder ResponseEncoderFunc,
	exec edge.Edge) *RequestHandler {

	rh := &RequestHandler{
		ctx,
		decoder,
		encoder,
		exec,
		false,
		nil,
		nil,
	}
	return rh
}

// Handles the request
func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := h.ctx

	// extract headers
	subdomain := strings.Split(r.Host, ".")
	fmt.Println("subdomain = ", subdomain[0])
	ctx = context.WithValue(ctx, "subdomain", subdomain[0])

	// Verify the access-token
	if h.protected {

		cookie, err := r.Cookie(AuthTokenCookieName)
		if err != nil {
			handleError(ctx, errUnauthorized, AUTH, w)
			return
		}

		if "" == cookie.Value {
			handleError(ctx, errUnauthorized, AUTH, w)
			return
		}

		token, err := auth.VefifyToken(h.verifier, cookie.Value)
		if err != nil || !token.Valid {
			fmt.Println(err)
			handleError(ctx, err, AUTH, w)
			return
		}

		ctx = context.WithValue(ctx, "token", auth.GetToken(token))

		// Refresh the token
		refreshtoken(token, h.signer, h.protected, w)
	}

	// decodes the request
	req, err := h.decode(ctx, r)
	if err != nil {
		handleError(ctx, err, DECODE, w)
		return
	}

	// process the request
	res, err := h.process(ctx, req)
	if err != nil {
		handleError(ctx, err, PROCESS, w)
		return
	}

	err = h.encode(ctx, w, res)
	if err != nil {
		handleError(ctx, err, ENCODE, w)
		return
	}
}

// HandleError handles the error by setting up the right message and status code
func handleError(ctx context.Context, err error, phase int, w http.ResponseWriter) {

	fmt.Println(err)
	switch phase {
	case AUTH:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case DECODE:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case PROCESS:
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	case ENCODE:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func refreshtoken(token *jwt.Token, signer interface{}, protected bool, w http.ResponseWriter) {

	if protected {
		signedTok, err := auth.RefreshToken(signer, token)
		if err != nil {
			fmt.Println("Failed to refresh the token.", err)
		}

		cookie := &http.Cookie{
			Name:     AuthTokenCookieName,
			Value:    signedTok,
			Expires:  time.Now().Add(time.Minute * 15),
			HttpOnly: true,
			Path:     "/",
			//Secure : true, // TODO enable this to ensure the cookie is passed only with https
		}
		http.SetCookie(w, cookie)
	}
}
