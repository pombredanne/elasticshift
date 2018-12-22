/*
Copyright 2017 The Elasticshift Authors.
*/
package providers

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
	"github.com/elasticshift/elasticshift/internal/shiftserver/sysconf"
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
	errNoProviderFound = "No provider found for %s : %v"
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

	Search(token, vcsName, repoName string) (types.Repository, error)

	CreateHook(token, owner, repo string) error
}

// Providers type
type Providers struct {
	logger *logrus.Entry
	loggr  logger.Loggr
	store  store.Sysconf
}

func New(loggr logger.Loggr, s store.Shift) Providers {
	return Providers{loggr: loggr, logger: loggr.GetLogger("oauth2/providers"), store: s.Sysconf}
}

// Get the provider by namee
func (p Providers) Get(name string) (Provider, error) {

	var conf types.VCSSysConf
	err := p.store.GetSysConf(sysconf.VcsKind, name, &conf)
	if err != nil {
		return nil, fmt.Errorf(errNoProviderFound, name, err)
	}

	fmt.Println("Providers.Get(): ", conf.Name)

	switch conf.Name {
	case GithubProviderName:
		return GithubProvider(p.loggr, conf.Key, conf.Secret, conf.CallbackURL, conf.HookURL), nil
	case GitlabProviderName:
		return GitlabProvider(p.loggr, conf.Key, conf.Secret, conf.CallbackURL, conf.HookURL), nil
	case BitbucketProviderName:
		return BitbucketProvider(p.loggr, conf.Key, conf.Secret, conf.CallbackURL, conf.HookURL), nil
	}

	return nil, fmt.Errorf("No provider found for ", name)
}
