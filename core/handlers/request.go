package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/palantir/stacktrace"
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
	Logger      *logrus.Logger
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
			err = stacktrace.Propagate(err, "Error occured during DECODE request")
			handleError(ctx, h.Logger, err, DECODE, w)
			return
		}
	}

	// process the request
	res, err := h.ProcessFunc(ctx, req)
	if err != nil {
		err = stacktrace.Propagate(err, "Error occured during PROCESS request")
		handleError(ctx, h.Logger, err, PROCESS, w)
		return
	}

	err = h.EncodeFunc(ctx, w, res)
	if err != nil {
		err = stacktrace.Propagate(err, "Error occured during ENCODE response")
		handleError(ctx, h.Logger, err, ENCODE, w)
		return
	}
}

// HandleError handles the error by setting up the right message and status code
func handleError(ctx context.Context, logger *logrus.Logger, err error, phase int, w http.ResponseWriter) {

	logger.Errorln(stacktrace.RootCause(err))

	switch phase {
	case DECODE:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case PROCESS:
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	case ENCODE:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}