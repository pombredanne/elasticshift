package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/conspico/esh/core/edge"
)

// Enum that represents the PHASE of the request
const (
	DECODE  = 1
	PROCESS = 2
	ENCODE  = 3
)

var (
	errNoProcessEdgeFound = errors.New("No edge found")
)

// RequestHandler ..
// Any request reaches to ESH server lands here
type RequestHandler struct {
	DecodeFunc  RequestDecoderFunc
	ProcessFunc edge.Edge
	EncodeFunc  ResponseEncoderFunc
}

// ServeHTTP ..
// Handle the request
func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// extract headers
	subdomain := strings.Split(r.Host, ".")
	ctx = context.WithValue(ctx, "subdomain", subdomain[0])

	// decodes the request
	var req interface{}
	var err error

	if h.DecodeFunc != nil {

		req, err = h.DecodeFunc(ctx, r)
		if err != nil {
			handleError(ctx, err, DECODE, w)
			return
		}
	}

	// process the request
	res, err := h.ProcessFunc(ctx, req)
	if err != nil {
		handleError(ctx, err, PROCESS, w)
		return
	}

	err = h.EncodeFunc(ctx, w, res)
	if err != nil {
		handleError(ctx, err, ENCODE, w)
		return
	}
}

// HandleError handles the error by setting up the right message and status code
func handleError(ctx context.Context, err error, phase int, w http.ResponseWriter) {

	fmt.Println(err)
	switch phase {
	case DECODE:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case PROCESS:
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	case ENCODE:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
