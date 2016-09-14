package http

import (
	"context"
	"net/http"

	"gitlab.com/conspico/esh/core/edge"
)

// Enum that represents the PHASE of the request
const (
	DECODE  = 0
	PROCESS = 1
	ENCODE  = 2
)

// RequestHandler ..
// Any request reaches to ESH server lands here
// and a life cycle will be performed such as decode, process, encode  a request
type RequestHandler struct {
	ctx     context.Context
	decode  RequestDecoderFunc
	encode  ResponseEncoderFunc
	process edge.Edge
}

// NewRequestHandler creates a reqeust handler for given edge
func NewRequestHandler(
	ctx context.Context,
	decoder RequestDecoderFunc,
	encoder ResponseEncoderFunc,
	exec edge.Edge) *RequestHandler {

	rh := &RequestHandler{
		ctx,
		decoder,
		encoder,
		exec,
	}
	return rh
}

// Handles the request
func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := h.ctx

	// extract headers and set it to ctx

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

	switch phase {
	case DECODE:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case PROCESS:
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	case ENCODE:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
