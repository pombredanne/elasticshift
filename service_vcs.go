// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"bytes"

	"encoding/base64"

	"github.com/Sirupsen/logrus"
	"github.com/palantir/stacktrace"
)

// expiryDelta determines how earlier a token should be considered
const expiryDelta = 10 * time.Second

// VCS account owner type
const (
	OwnerTypeUser = "user"
	OwnerTypeOrg  = "org"
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
	teamDS       TeamDatastore
	repoDS       RepoDatastore
	sysconfDS    SysconfDatastore
	vcsProviders *Providers
	config       Config
	logger       *logrus.Logger
	vcsConf      map[string]VCSSysConf
}

// NewVCSService ..
func NewVCSService(ctx AppContext) VCSService {

	conf := ctx.Config
	/*providers := NewProviders(
		GithubProvider(ctx.Logger, conf.Github.Key, conf.Github.Secret, conf.Github.Callback),
		GitlabProvider(ctx.Logger, conf.Gitlab.Key, conf.Gitlab.Secret, conf.Gitlab.Callback),
		BitbucketProvider(ctx.Logger, conf.Bitbucket.Key, conf.Bitbucket.Secret, conf.Bitbucket.Callback),
	)*/

	return &vcsService{
		teamDS:    ctx.TeamDatastore,
		repoDS:    ctx.RepoDatastore,
		sysconfDS: ctx.SysconfDatastore,
		config:    conf,
		logger:    ctx.Logger,
		vcsConf:   make(map[string]VCSSysConf),
	}
}

func (s vcsService) Authorize(teamname, provider string, r *http.Request) (AuthorizeResponse, error) {

	p, err := s.getProvider(provider)
	if err != nil {
		return AuthorizeResponse{}, stacktrace.Propagate(err, "Getting provider %s failed", provider)
	}

	// Get the base URL
	var buf bytes.Buffer
	buf.WriteString(teamname)
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

	p, err := s.getProvider(provider)
	if err != nil {
		return AuthorizeResponse{}, stacktrace.Propagate(err, "Getting provider %s failed", provider)
	}

	redirectURL := p.GetRedirectURL(id)
	v, err := p.Authorized(code, redirectURL)
	if err != nil {
		return AuthorizeResponse{}, stacktrace.Propagate(err, "Finalize the authorization failed.")
	}

	unescID := s.decode(id)
	escID := strings.Split(unescID, SEMICOLON)
	fmt.Println(escID)

	// persist user
	tname := escID[0]

	acc, err := s.teamDS.GetVCSByID(tname, v.ID)
	if strings.EqualFold(acc.ID, v.ID) {

		acc.AccessToken = v.AccessToken
		acc.AccessCode = v.AccessCode
		acc.RefreshToken = v.RefreshToken
		acc.OwnerType = v.OwnerType
		acc.TokenExpiry = v.TokenExpiry

		s.teamDS.UpdateVCS(tname, acc)

		return AuthorizeResponse{Conflict: true, Err: errVCSAccountAlreadyLinked, Request: r}, nil
	}

	err = s.teamDS.SaveVCS(tname, &v)

	url := escID[1] + "/api/vcs"
	return AuthorizeResponse{Err: nil, URL: url, Request: r}, err
}

func (s vcsService) GetVCS(team string) (GetVCSResponse, error) {

	result, err := s.teamDS.GetTeam(team)
	return GetVCSResponse{Result: result.Accounts}, err
}

func (s vcsService) SyncVCS(team, userName, id string) (bool, error) {

	acc, err := s.teamDS.GetVCSByID(team, id)
	if err != nil {
		return false, stacktrace.Propagate(err, "Get by VCS ID failed during sync.")
	}

	err = s.sync(acc, userName, team)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s vcsService) sync(acc VCS, userName, team string) error {

	// Get the token
	t, err := s.getToken(team, acc)
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
	lrpo, err := s.repoDS.GetReposByVCSID(team, acc.ID)
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

			updated := false
			if r.Name != rp.Name {
				r.Name = rp.Name
				updated = true
			}

			if r.Private != rp.Private {
				r.Private = rp.Private
				updated = true
			}

			if r.Link != rp.Link {
				r.Link = rp.Link
				updated = true
			}

			if r.Description != rp.Description {
				r.Description = rp.Description
				updated = true
			}

			if r.Fork != rp.Fork {
				r.Fork = rp.Fork
				updated = true
			}

			if r.DefaultBranch != rp.DefaultBranch {
				r.DefaultBranch = rp.DefaultBranch
				updated = true
			}

			if r.Language != rp.Language {
				r.Language = rp.Language
				updated = true
			}

			if updated {
				// perform update
				s.repoDS.Update(r)
			}
		} else {

			// perform insert
			rp.Team = team
			rp.VcsID = acc.ID
			err := s.repoDS.Save(&rp)
			if err != nil {
				s.logger.Errorln(err)
			}
		}

		// removes from the map
		if exist {
			delete(rpo, r.RepoID)
		}
	}

	var ids []bson.ObjectId
	// Now iterate thru deleted repositories.
	for _, rp := range rpo {
		ids = append(ids, rp.ID)
	}

	if len(ids) > 0 {
		err = s.repoDS.DeleteIds(ids)
		if err != nil {
			return stacktrace.Propagate(err, "Failed to delete the vcs that does not exist remotly")
		}
	}

	return nil
}

// Gets the valid token
// Checks whether the token is expired.
// Expired token will get refreshed.
func (s vcsService) getToken(team string, a VCS) (string, error) {

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
	a.AccessToken = tok.AccessToken
	a.TokenExpiry = tok.Expiry
	a.RefreshToken = tok.RefreshToken
	a.TokenType = tok.TokenType
	err = s.teamDS.UpdateVCS(team, a)

	if err != nil {
		return "", stacktrace.Propagate(err, "Failed to update VCS after token refreshed.")
	}
	return tok.AccessToken, nil
}

// Gets the provider by type
func (s vcsService) getProvider(providerName string) (Provider, error) {

	if len(s.vcsConf) == 0 {

		conf, err := s.sysconfDS.GetVCSTypes()
		if err != nil {
			return nil, err
		}

		for _, c := range conf {
			(s.vcsConf)[c.Name] = c
		}
	}

	v := s.vcsConf[providerName]
	var p Provider
	if strings.EqualFold(GithubProviderName, v.Name) {
		p = GithubProvider(s.logger, v.Key, v.Secret, v.CallbackURL)
	} else if strings.EqualFold(BitBucketProviderName, v.Name) {
		p = BitbucketProvider(s.logger, v.Key, v.Secret, v.CallbackURL)
	} else if strings.EqualFold(GitlabProviderName, v.Name) {
		p = GitlabProvider(s.logger, v.Key, v.Secret, v.CallbackURL)
	}
	return p, nil
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
