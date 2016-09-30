package vcs

import (
	"context"

	"gitlab.com/conspico/esh/core/edge"
)

func makeAuthorizeEdge(s Service) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthorizeRequest)
		return s.Authorize(req.Domain, req.Provider, req.Request)
	}
}
