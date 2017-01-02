// Package esh ...
// Author: Ghazni Nattarshah
// Date: Oct 01, 2016
package esh

import (
	"context"

	"gitlab.com/conspico/esh/core/edge"
)

/***** TEAM ****/
func makeCreateTeamEdge(s TeamService) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createTeamRequest)
		created, err := s.Create(req.Name)
		return createTeamResponse{Created: created, Err: err}, nil
	}
}

/***** USER ****/
func makeSignupEdge(s UserService) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		token, err := s.Create(request.(signupRequest))
		return signInResponse{Token: token, Err: err}, nil
	}
}

func makeSignInEdge(s UserService) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		token, err := s.SignIn(request.(signInRequest))
		return signInResponse{Token: token, Err: err}, nil
	}
}

func makeSignOutEdge(s UserService) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(signOut)
		s.SignOut()
		return req, nil
	}
}

func makeVerifyCodeEdge(s UserService) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(verifyCodeRequest)
		valid, err := s.Verify(req.Code)
		return genericResponse{Valid: valid, Err: err}, nil
	}
}

/**** VCS ****/
func makeAuthorizeEdge(s VCSService) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return s.Authorize(request.(AuthorizeRequest))
	}
}

func makeAuthorizedEdge(s VCSService) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return s.Authorized(request.(AuthorizeRequest))
	}
}

func makeGetVCSEdge(s VCSService) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		teamID := request.(string)
		return s.GetVCS(teamID)
	}
}

func makeSyncVCSEdge(s VCSService) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return s.SyncVCS(request.(SyncVCSRequest))
	}
}

/**** REPO ****/
func makeGetRepoEdge(s RepoService) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRepoRequest)
		if req.VCSID == "" {
			return s.GetRepos(req.TeamID)
		}
		return s.GetReposByVCSID(req.TeamID, req.VCSID)
	}
}
