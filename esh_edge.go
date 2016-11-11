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
		req := request.(signupRequest)
		token, err := s.Create(req.Team, req.Domain, req.Fullname, req.Email, req.Password)
		return signInResponse{Token: token, Err: err}, nil
	}
}

func makeSignInEdge(s UserService) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(signInRequest)
		token, err := s.SignIn(req.Team, req.Domain, req.Email, req.Password)
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
		req := request.(AuthorizeRequest)
		return s.Authorize(req.TeamID, req.Provider, req.Request)
	}
}

func makeAuthorizedEdge(s VCSService) edge.Edge {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuthorizeRequest)
		return s.Authorized(req.ID, req.Provider, req.Code, req.Request)
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
		req := request.(SyncVCSRequest)
		return s.SyncVCS(req.TeamID, req.Username, req.ProviderID)
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
