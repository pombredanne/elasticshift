package vcs

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/conspico/esh/core/util"
	"gitlab.com/conspico/esh/team"

	"github.com/spf13/viper"
)

// expiryDelta determines how earlier a token should be considered
const expiryDelta = 10 * time.Second

// Service ..
type Service interface {
	Authorize(subdomain, provider string, r *http.Request) (AuthorizeResponse, error)
	Authorized(subdomain, provider, code string, r *http.Request) (AuthorizeResponse, error)
	GetVCS(teamID string) (GetVCSResponse, error)
	SyncVCS(teamID, provider string) (bool, error)
}

type service struct {
	vcsDS        Datastore
	teamDS       team.Datastore
	vcsProviders *Providers
	config       *viper.Viper
}

// NewService ..
func NewService(v Datastore, t team.Datastore, conf *viper.Viper) Service {

	providers := NewProviders(
		GithubProvider(conf.GetString("github.key"), conf.GetString("github.secret"), conf.GetString("github.callback")),
		BitbucketProvider(conf.GetString("bitbucket.key"), conf.GetString("bitbucket.secret"), conf.GetString("bitbucket.callback")),
	)

	return &service{
		vcsProviders: providers,
		vcsDS:        v,
		teamDS:       t,
		config:       conf,
	}
}

func (s service) Authorize(teamID, provider string, r *http.Request) (AuthorizeResponse, error) {

	p, err := s.vcsProviders.Get(provider)
	if err != nil {
		return AuthorizeResponse{}, err
	}

	url := p.Authorize(teamID)

	return AuthorizeResponse{Err: nil, URL: url, Request: r}, nil
}

// Authorized ..
// Invoked when authorization finished by oauth app
func (s service) Authorized(teamID, provider, code string, r *http.Request) (AuthorizeResponse, error) {

	p, err := s.vcsProviders.Get(provider)
	if err != nil {
		return AuthorizeResponse{}, err
	}

	u, err := p.Authorized(code)
	if err != nil {
		fmt.Println(err)
		return AuthorizeResponse{}, err
	}

	// persist user
	u.ID, _ = util.NewUUID()
	u.TeamID = teamID
	u.CreatedDt = time.Now()
	u.UpdatedDt = time.Now()
	err = s.vcsDS.Save(&u)

	url := "/api/vcs"
	return AuthorizeResponse{Err: nil, URL: url, Request: r}, err
}

func (s service) GetVCS(teamID string) (GetVCSResponse, error) {

	result, err := s.vcsDS.GetVCS(teamID)
	return GetVCSResponse{Result: result}, err
}

func (s service) SyncVCS(teamID, providerID string) (bool, error) {

	acc, err := s.vcsDS.GetByID(providerID)
	if err != nil {
		return false, err
	}

	err = s.sync(acc)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s service) sync(acc VCS) error {

	// Get the token
	/*t, err := s.getToken(acc)
	if err != nil {
		return fmt.Errorf(errGetUpdatedFokenFailed, err)
	}

	// fetch the existing repository
	p, err := s.getProvider(acc.Type)
	if err != nil {
		return fmt.Errorf(errNoProviderFound, err)
	}

	// repository received from provider
	repos, err := p.GetRepos(t, acc.OwnerType) */

	// combine the result set

	// insert or update the repository

	return nil
}

// Gets the valid token
// Checks whether the token is expired.
// Expired token will get refreshed.
func (s service) getToken(a VCS) (string, error) {

	// Never expire type token
	if a.RefreshToken == "" {
		return a.AccessToken, nil
	}

	// Token that requires frequent refresh
	// check if the token is expired
	if !a.TokenExpiry.Add(-expiryDelta).Before(time.Now()) {
		return a.AccessToken, nil
	}

	p, err := s.getProvider(a.Type)
	if err != nil {
		return "", err
	}

	// Refresh the token
	tok, err := p.RefreshToken(a.RefreshToken)

	// persist the updated token information
	err = s.vcsDS.Update(&a, VCS{
		AccessToken:  tok.AccessToken,
		TokenExpiry:  tok.Expiry,
		RefreshToken: tok.RefreshToken,
		TokenType:    tok.TokenType,
	})

	if err != nil {
		return "", err
	}
	return tok.AccessToken, nil
}

// Gets the provider by type
func (s service) getProvider(vcsType int) (Provider, error) {

	var name string
	switch vcsType {
	case GithubType:
		name = GithubProviderName
	case BitBucketType:
		name = BitBucketProviderName
	}

	return s.vcsProviders.Get(name)
}
