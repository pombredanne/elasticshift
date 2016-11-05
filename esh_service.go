package esh

import "net/http"

// TeamService ..
type TeamService interface {
	Create(name string) (bool, error)
}

// UserService ..
type UserService interface {
	Create(teamName, domain, fullName, email, password string) (string, error)
	SignIn(teamName, domain, email, password string) (string, error)
	SignOut() (bool, error)
	Verify(code string) (bool, error)
}

// VCSService ..
type VCSService interface {
	Authorize(subdomain, provider string, r *http.Request) (AuthorizeResponse, error)
	Authorized(subdomain, provider, code string, r *http.Request) (AuthorizeResponse, error)
	GetVCS(teamID string) (GetVCSResponse, error)
	SyncVCS(teamID, userName, provider string) (bool, error)
}
