package handlers

import (
	"context"
	"net/http"
)

// ResponseDecoderFunc is used after received a response from a call to an edge (http endpoint)
// and parsing the response data.
// Ex: A call to an edge (htpt endpoint) and parsing a json response
type ResponseDecoderFunc func(context.Context, *http.Response) (interface{}, error)

// ResponseEncoderFunc is used after received a response from a call to an edge (http endpoint)
// and parsing the response data
// Ex: Encoding a json response will then after returned to the caller
type ResponseEncoderFunc func(context.Context, http.ResponseWriter, interface{}) error

// ErrorHandlerFunc is used to identify the type of errors such as bad request,
// internal server error. It also writes the error message to user and
// set the staus code in header
type ErrorHandlerFunc func(context.Context, error, int, http.ResponseWriter)
