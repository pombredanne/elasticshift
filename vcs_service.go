package esh

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/conspico/esh/core/util"
)

// expiryDelta determines how earlier a token should be considered
const expiryDelta = 10 * time.Second

// VCS account owner type
const (
	OwnerTypeUser = 1
	OwnerTypeOrg  = 2
)

type vcsService struct {
	vcsDS        VCSDatastore
	teamDS       TeamDatastore
	vcsProviders *Providers
	config       Config
}

// NewVCSService ..
func NewVCSService(v VCSDatastore, t TeamDatastore, conf Config) VCSService {

	providers := NewProviders(
		GithubProvider(conf.Github.Key, conf.Github.Secret, conf.Github.Callback),
		BitbucketProvider(conf.Bitbucket.Key, conf.Bitbucket.Secret, conf.Bitbucket.Callback),
	)

	return &vcsService{
		vcsProviders: providers,
		vcsDS:        v,
		teamDS:       t,
		config:       conf,
	}
}

func (s vcsService) Authorize(teamID, provider string, r *http.Request) (AuthorizeResponse, error) {

	p, err := s.vcsProviders.Get(provider)
	if err != nil {
		return AuthorizeResponse{}, err
	}

	url := p.Authorize(teamID)

	return AuthorizeResponse{Err: nil, URL: url, Request: r}, nil
}

// Authorized ..
// Invoked when authorization finished by oauth app
func (s vcsService) Authorized(teamID, provider, code string, r *http.Request) (AuthorizeResponse, error) {

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

func (s vcsService) GetVCS(teamID string) (GetVCSResponse, error) {

	result, err := s.vcsDS.GetVCS(teamID)
	return GetVCSResponse{Result: result}, err
}

func (s vcsService) SyncVCS(teamID, providerID string) (bool, error) {

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

func (s vcsService) sync(acc VCS) error {

	// Get the token
	t, err := s.getToken(acc)
	if err != nil {
		return fmt.Errorf(errGetUpdatedFokenFailed, err)
	}

	// fetch the existing repository
	p, err := s.getProvider(acc.Type)
	if err != nil {
		return fmt.Errorf(errNoProviderFound, err)
	}

	// repository received from provider
	repos, err := p.GetRepos(t, acc.OwnerType)
	if err != nil {
		return err
	}
	fmt.Println(repos)

	// Fetch the repositories from VCS
	/*lrpo, err := s.repoDS.GetReposByVCSID(acc.ID)
	if err != nil {
		return err
	}
	rpo := make(map[string]repo.Repo)
	for _, l := range lrpo {
		rpo[l.RepoID] = l
	}*/

	// combine the result set

	// insert or update the repository

	return nil
}

// Gets the valid token
// Checks whether the token is expired.
// Expired token will get refreshed.
func (s vcsService) getToken(a VCS) (string, error) {

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
	err = s.vcsDS.UpdateVCS(&a, VCS{
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
func (s vcsService) getProvider(vcsType int) (Provider, error) {

	var name string
	switch vcsType {
	case GithubType:
		name = GithubProviderName
	case BitBucketType:
		name = BitBucketProviderName
	}

	return s.vcsProviders.Get(name)
}
