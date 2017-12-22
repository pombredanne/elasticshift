/*
Copyright 2017 The Elasticshift Authors.
*/
package providers

import (
	"fmt"
	"time"

	"gitlab.com/conspico/elasticshift/api/types"
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

// Account owner type
const (
	OwnerTypeUser = "user"
	OwnerTypeOrg  = "org"
)

// True or False
const (
	True  = 1
	False = 0
)

var (
	errNoProviderFound = "No provider found :"
)

// Token ..
type Token struct {
	AccessToken string `json:"access_token"`

	// (bearer, mac, etc)
	TokenType string `json:"token_type"`

	// The refresh token, which can be used to obtain new
	// access tokens using the same authorization grant
	RefreshToken string `json:"refresh_token"`

	// The lifetime in seconds of the access token.
	ExpiresIn int64 `json:"expires_in"`

	Expiry time.Time `json:"expiry,omitempty"`

	CreatedAt int64 `json:"created_at"`
	Scope     string
}

// Provider ..
type Provider interface {
	Name() string

	Authorize(baseURL string) string

	Authorized(id, code string) (types.VCS, error)

	RefreshToken(token string) (*oauth2.Token, error)

	GetRepos(token, accountName string, owner string) ([]types.Repository, error)

	CreateHook(token, owner, repo string) error
}

// Providers type
type Providers struct {
	Providers map[string]Provider
}

func New() Providers {
	return Providers{make(map[string]Provider)}
}

// NewProviders ...
func NewProviders(pvider ...Provider) Providers {

	var prov = make(map[string]Provider)
	for _, p := range pvider {
		prov[p.Name()] = p
	}
	return Providers{prov}
}

// Set the provider for the given name
func (prov Providers) Set(name string, p Provider) {
	prov.Providers[name] = p
}

// Get the provider by namee
func (prov Providers) Get(name string) (Provider, error) {

	if p, ok := prov.Providers[name]; ok {
		return p, nil
	}

	return nil, fmt.Errorf(errNoProviderFound, name)
}
