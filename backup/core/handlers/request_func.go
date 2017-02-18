package handlers

import (
	"context"
	"net/http"
)

// RequestFunc parse necessary http header common and place it in context
// Ex: header information, create request trace id
type RequestFunc func(context.Context, *http.Request) context.Context

// RequestDecoderFunc is used when receiving a call to an edge (http endpoint)
// and parsing the request data
// Ex: A call to an edge (htpt endpoint) and parsing a json request for further process
type RequestDecoderFunc func(context.Context, *http.Request) (interface{}, error)

// RequestEncoderFunc is used when sendina a call to an edge (http endpoint)
// and constructing a request data
// Ex: A call to an edge (htpt endpoint) and encoding a json request
type RequestEncoderFunc func(context.Context, *http.Request, interface{}) error
