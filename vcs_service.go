package esh

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"bytes"

	"encoding/base64"

	"github.com/Sirupsen/logrus"
	"github.com/palantir/stacktrace"
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

// Constants for performing encode decode
const (
	EQUAL        = "="
	DOUBLEEQUALS = "=="
	DOT0         = ".0"
	DOT1         = ".1"
	DOT2         = ".2"
)

type vcsService struct {
	vcsDS        VCSDatastore
	teamDS       TeamDatastore
	repoDS       RepoDatastore
	vcsProviders *Providers
	config       Config
	logger       *logrus.Logger
}

// NewVCSService ..
func NewVCSService(ctx AppContext) VCSService {

	conf := ctx.Config
	providers := NewProviders(
		GithubProvider(ctx.Logger, conf.Github.Key, conf.Github.Secret, conf.Github.Callback),
		BitbucketProvider(ctx.Logger, conf.Bitbucket.Key, conf.Bitbucket.Secret, conf.Bitbucket.Callback),
	)

	return &vcsService{
		vcsProviders: providers,
		vcsDS:        ctx.VCSDatastore,
		teamDS:       ctx.TeamDatastore,
		repoDS:       ctx.RepoDatastore,
		config:       conf,
		logger:       ctx.Logger,
	}
}

func (s vcsService) Authorize(teamID, provider string, r *http.Request) (AuthorizeResponse, error) {

	p, err := s.vcsProviders.Get(provider)
	if err != nil {
		return AuthorizeResponse{}, stacktrace.Propagate(err, "Getting provider %s failed", provider)
	}

	// Get the base URL
	var buf bytes.Buffer
	buf.WriteString(teamID)
	buf.WriteString(SEMICOLON)
	buf.WriteString(SLASH)
	buf.WriteString(SLASH)
	buf.WriteString(r.Host)

	url := p.Authorize(s.encode(buf.String()))

	return AuthorizeResponse{Err: nil, URL: url, Request: r}, nil
}

// Authorized ..
// Invoked when authorization finished by oauth app
func (s vcsService) Authorized(id, provider, code string, r *http.Request) (AuthorizeResponse, error) {

	p, err := s.vcsProviders.Get(provider)
	if err != nil {
		return AuthorizeResponse{}, stacktrace.Propagate(err, "Getting provider %s failed", provider)
	}

	u, err := p.Authorized(code)
	if err != nil {
		return AuthorizeResponse{}, stacktrace.Propagate(err, "Finalize the authorization failed.")
	}

	unescID := s.decode(id)
	escID := strings.Split(unescID, SEMICOLON)

	// persist user
	u.TeamID = escID[0]

	acc, err := s.vcsDS.GetByVCSID(u.TeamID, u.VcsID)
	if strings.EqualFold(acc.VcsID, u.VcsID) {

		updvcs := VCS{}
		updvcs.UpdatedDt = time.Now()
		updvcs.AccessToken = u.AccessToken
		updvcs.AccessCode = u.AccessCode
		updvcs.RefreshToken = u.RefreshToken
		updvcs.OwnerType = u.OwnerType
		updvcs.TokenExpiry = u.TokenExpiry

		s.vcsDS.Update(&acc, updvcs)

		s.logger.Info("Conflict")
		return AuthorizeResponse{Conflict: true, Err: errVCSAccountAlreadyLinked, Request: r}, nil
	}

	u.ID, _ = util.NewUUID()
	u.CreatedDt = time.Now()
	u.UpdatedDt = time.Now()

	err = s.vcsDS.Save(&u)

	url := escID[1] + "/api/vcs"
	return AuthorizeResponse{Err: nil, URL: url, Request: r}, err
}

func (s vcsService) GetVCS(teamID string) (GetVCSResponse, error) {

	result, err := s.vcsDS.GetVCS(teamID)
	return GetVCSResponse{Result: result}, err
}

func (s vcsService) SyncVCS(teamID, userName, providerID string) (bool, error) {

	acc, err := s.vcsDS.GetByID(providerID)
	if err != nil {
		return false, stacktrace.Propagate(err, "Get by VCS ID failed during sync.")
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
		return stacktrace.Propagate(fmt.Errorf(errGetUpdatedFokenFailed, err), "Get token failed")
	}

	// fetch the existing repository
	p, err := s.getProvider(acc.Type)
	if err != nil {
		return stacktrace.Propagate(fmt.Errorf(errNoProviderFound, err), errNoProviderFound)
	}

	// repository received from provider
	repos, err := p.GetRepos(t, acc.Name, acc.OwnerType)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to get repos from provider %s", p.Name())
	}

	// Fetch the repositories from esh repo store
	lrpo, err := s.repoDS.GetReposByVCSID(acc.TeamID, acc.ID)
	if err != nil {
		return stacktrace.Propagate(err, "Getting repos by vcs id failed.")
	}

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

	var ids []string
	// Now iterate thru deleted repositories.
	for _, rp := range rpo {
		ids = append(ids, rp.ID)
	}

	err = s.repoDS.DeleteIds(ids)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to delete the vcs that does not exist remotly")
	}

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
		return "", stacktrace.Propagate(fmt.Errorf(errNoProviderFound, err), errNoProviderFound)
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
		return "", stacktrace.Propagate(err, "Failed to update VCS after token refreshed.")
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

func (s vcsService) encode(id string) string {

	eid := base64.URLEncoding.EncodeToString([]byte(id))
	if strings.Contains(eid, DOUBLEEQUALS) {
		eid = strings.TrimRight(eid, DOUBLEEQUALS) + DOT2
	} else if strings.Contains(eid, EQUAL) {
		eid = strings.TrimRight(eid, EQUAL) + DOT1
	} else {
		eid = eid + DOT0
	}
	return eid
}

func (s vcsService) decode(id string) string {

	if strings.Contains(id, DOT2) {
		id = strings.TrimRight(id, DOT2) + DOUBLEEQUALS
	} else if strings.Contains(id, DOT1) {
		id = strings.TrimRight(id, DOT1) + EQUAL
	} else {
		id = strings.TrimRight(id, DOT0)
	}
	did, _ := base64.URLEncoding.DecodeString(id)
	return string(did[:])
}
