package esh

import (
	"fmt"

	"golang.org/x/oauth2"
)

// VCS user type
const (
	GithubType    = 1
	GitlabType    = 2
	BitBucketType = 3
	SvnType       = 4
	TfsType       = 5
)

// Provider ..
type Provider interface {
	Name() string

	Authorize(team string) string

	Authorized(code string) (VCS, error)

	RefreshToken(token string) (*oauth2.Token, error)

	GetRepos(token string, owner int) ([]Repo, error)
}

// Providers type
type Providers struct {
	Providers map[string]Provider
}

// NewProviders ...
func NewProviders(pvider ...Provider) *Providers {

	var prov = make(map[string]Provider)
	for _, p := range pvider {
		prov[p.Name()] = p
	}
	return &Providers{prov}
}

// Get the provider by namee
func (prov Providers) Get(name string) (Provider, error) {

	p := prov.Providers[name]
	if p == nil {
		return nil, fmt.Errorf(errNoProviderFound, name)
	}
	return p, nil
}
