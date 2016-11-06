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

// True or False
const (
	True  = 1
	False = 0
)

type vcsService struct {
	vcsDS        VCSDatastore
	teamDS       TeamDatastore
	repoDS       RepoDatastore
	vcsProviders *Providers
	config       Config
}

// NewVCSService ..
func NewVCSService(v VCSDatastore, t TeamDatastore, r RepoDatastore, conf Config) VCSService {

	providers := NewProviders(
		GithubProvider(conf.Github.Key, conf.Github.Secret, conf.Github.Callback),
		BitbucketProvider(conf.Bitbucket.Key, conf.Bitbucket.Secret, conf.Bitbucket.Callback),
	)

	return &vcsService{
		vcsProviders: providers,
		vcsDS:        v,
		teamDS:       t,
		repoDS:       r,
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

func (s vcsService) SyncVCS(teamID, userName, providerID string) (bool, error) {

	acc, err := s.vcsDS.GetByID(providerID)
	if err != nil {
		return false, err
	}

	err = s.sync(acc, userName)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s vcsService) sync(acc VCS, userName string) error {

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
	repos, err := p.GetRepos(t, acc.Name, acc.OwnerType)
	if err != nil {
		return err
	}

	// Fetch the repositories from esh repo store
	lrpo, err := s.repoDS.GetReposByVCSID(acc.ID)
	if err != nil {
		return err
	}
	fmt.Println("Local repo count = ", len(lrpo))
	rpo := make(map[string]Repo)
	for _, l := range lrpo {
		rpo[l.RepoID] = l
	}

	// combine the result set
	for _, rp := range repos {

		r, exist := rpo[rp.RepoID]
		if exist {

			updrepo := Repo{}
			updated := false
			if r.Name != rp.Name {
				updrepo.Name = rp.Name
				updated = true
			}

			if r.Private != rp.Private {
				updrepo.Private = rp.Private
				updated = true
			}

			if r.Link != rp.Link {
				updrepo.Link = rp.Link
				updated = true
			}

			if r.Description != rp.Description {
				updrepo.Description = rp.Description
				updated = true
			}

			if r.Fork != rp.Fork {
				updrepo.Fork = rp.Fork
				updated = true
			}

			if r.DefaultBranch != rp.DefaultBranch {
				updrepo.DefaultBranch = rp.DefaultBranch
				updated = true
			}

			if r.Language != rp.Language {
				updrepo.Language = rp.Language
				updated = true
			}

			if updated {
				// perform update
				updrepo.UpdatedBy = userName
				s.repoDS.Update(r, updrepo)
			}
		} else {

			// perform insert
			rp.ID, _ = util.NewUUID()
			rp.CreatedDt = time.Now()
			rp.UpdatedDt = time.Now()
			rp.CreatedBy = userName
			rp.TeamID = acc.TeamID
			rp.VcsID = acc.ID
			s.repoDS.Save(&rp)
		}

		// removes from the map
		if exist {
			delete(rpo, r.RepoID)
		}
	}

	// Now iterate thru deleted repositories.
	/*for _, rp := range rpo {
		s.repoDS.Delete(rp)
	}*/

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
		fmt.Println(err)
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