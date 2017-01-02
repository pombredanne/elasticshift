// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

// TeamService ..
type TeamService interface {
	Create(name string) (bool, error)
}

// UserService ..
type UserService interface {
	Create(r signupRequest) (string, error)
	SignIn(r signInRequest) (string, error)
	SignOut() (bool, error)
	Verify(code string) (bool, error)
}

// VCSService ..
type VCSService interface {
	Authorize(r AuthorizeRequest) (AuthorizeResponse, error)
	Authorized(r AuthorizeRequest) (AuthorizeResponse, error)
	GetVCS(teamID string) (GetVCSResponse, error)
	SyncVCS(r SyncVCSRequest) (bool, error)
}

// RepoService ..
type RepoService interface {
	GetRepos(team string) (GetRepoResponse, error)
	GetReposByVCSID(team, id string) (GetRepoResponse, error)
}
