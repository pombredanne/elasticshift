package vcs

import (
	"context"

	"gitlab.com/conspico/esh/core/edge"
)

func makeAuthorizeEdge(s Service) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthorizeRequest)
		return s.Authorize(req.TeamID, req.Provider, req.Request)
	}
}

func makeAuthorizedEdge(s Service) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthorizeRequest)
		return s.Authorized(req.TeamID, req.Provider, req.Code, req.Request)
	}
}

func makeGetVCSEdge(s Service) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		teamID := request.(string)
		return s.GetVCS(teamID)
	}
}
