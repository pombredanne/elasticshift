/*
Copyright 2018 The Elasticshift Authors.
*/
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	"github.com/elasticshift/elasticshift/internal/shiftserver/identity/oauth2/providers"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
	"github.com/elasticshift/elasticshift/internal/shiftserver/team"
	"gopkg.in/mgo.v2/bson"
)

// expiryDelta determines how earlier a token should be considered
const expiryDelta = 10 * time.Second

var (
	// repository errors
	errNoURIProvided = errors.New("URI is empty")
)

// Resolver ...
type Resolver interface {
	FetchRepository(params graphql.ResolveParams) (interface{}, error)
	FetchBuild(params graphql.ResolveParams) (interface{}, error)
	AddRepository(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	store      store.Repository
	teamStore  store.Team
	buildStore store.Build
	logger     *logrus.Entry
	providers  providers.Providers
}

// NewResolver ...
func NewResolver(ctx context.Context, loggr logger.Loggr, s store.Shift) (Resolver, error) {

	r := &resolver{
		store:      s.Repository,
		teamStore:  s.Team,
		buildStore: s.Build,
		logger:     loggr.GetLogger("graphql/repository"),
	}
	return r, nil
}

func (r resolver) FetchRepository(params graphql.ResolveParams) (interface{}, error) {

	teamName, _ := params.Args["team"].(string)
	if teamName == "" {
		return nil, team.ErrTeamNameIsEmpty
	}

	vcsID, _ := params.Args["vcs_id"].(string)

	result, err := r.store.GetRepository(teamName, vcsID)

	var res types.RepositoryList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}

func (r resolver) FetchBuild(params graphql.ResolveParams) (interface{}, error) {

	id := params.Source.(types.Repository).ID.Hex()
	result, err := r.buildStore.FetchBuildByRepositoryID(id)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return nil, err
	}

	var res types.BuildList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}

func (r resolver) AddRepository(params graphql.ResolveParams) (interface{}, error) {

	uri, _ := params.Args["uri"].(string)
	teamName, _ := params.Args["team"].(string)

	if uri == "" {
		return nil, errNoURIProvided
	}

	if teamName == "" {
		return nil, team.ErrTeamNameIsEmpty
	}

	// parse uri and identify the VCS
	// git@github.com:nshahm/hybrid.test.runner.git
	protoGit := strings.HasPrefix(uri, "git@")
	protoHttps := strings.HasPrefix(uri, "http")

	eIdx := strings.LastIndex(uri, "/")
	var sIdx int
	var source, vcsName, repoName string
	if protoGit {

		sIdx = strings.Index(uri, "@")
		val := uri[sIdx+1 : eIdx]
		valArr := strings.Split(val, ":")
		source = valArr[0]
		vcsName = valArr[1]
		repoName = uri[eIdx+1:]

	} else if protoHttps {

		valArr := strings.Split(uri, "/")
		source = valArr[2]
		vcsName = valArr[3]
		repoName = valArr[4]
	}

	// TODO add source to repository such as github.com

	// parse the repository name
	account, err := r.teamStore.GetVCSByName(teamName, vcsName, source)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the vcs account %s for %s and team %s: %v", vcsName, source, teamName, err)
	}

	if account == nil {
		return nil, fmt.Errorf("Account '%s' from %s has not been linked with team '%s'", vcsName, source, teamName)
	}

	p, err := r.providers.Get(account.Kind)
	if err != nil {
		return nil, fmt.Errorf("No account named %s from %s linked: %v", vcsName, source, err)
	}

	token, err := r.getToken(teamName, *account, p)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch/refresh the token for account %s from %s: %v", vcsName, source, err)
	}

	// fetch the repository from VCSServer
	repoName = strings.TrimSuffix(repoName, ".git")

	repo, err := p.Search(token, vcsName, repoName)
	if err != nil {
		return nil, fmt.Errorf("Repository fetch from %s for %s failed: %v", source, vcsName, err)
	}

	if repo.RepoID == "" {
		return nil, fmt.Errorf("Project/repository not found for URI '%s'", uri)
	}

	repo.Team = teamName
	repo.VcsID = account.ID
	repo.Identifier = strings.Join([]string{source, vcsName}, "/")

	var currentRepo types.Repository
	err = r.store.FindOne(bson.M{"repo_id": repo.RepoID, "team": teamName, "name": repoName, "identifier": repo.Identifier}, &currentRepo)
	if err != nil && err.Error() != "not found" {
		return nil, fmt.Errorf("Failed to check the existance of the repository :%v", err)
	}

	if strings.EqualFold(currentRepo.RepoID, repo.RepoID) {
		return nil, fmt.Errorf("URI '%s' already added as a repository to your team", uri)
	}

	// Store the repository, if it doesn't exist, otherwise throw error
	err = r.store.Save(&repo)

	return repo, err
}

// Gets the valid token
// Checks whether the token is expired.
// Expired token will get refreshed.
func (r resolver) getToken(team string, a types.VCS, p providers.Provider) (string, error) {

	// Never expire type token
	if a.RefreshToken == "" {
		return a.AccessToken, nil
	}

	// Token that requires frequent refresh
	// check if the token is expired
	if !a.TokenExpiry.Add(-expiryDelta).Before(time.Now()) {
		return a.AccessToken, nil
	}

	// Refresh the token
	tok, err := p.RefreshToken(a.RefreshToken)

	a.AccessToken = tok.AccessToken
	a.TokenExpiry = tok.Expiry
	a.RefreshToken = tok.RefreshToken

	// persist the updated token information
	err = r.teamStore.UpdateVCS(team, a)

	if err != nil {
		return "", fmt.Errorf("Failed to update VCS after token refreshed: %v", err)
	}
	return tok.AccessToken, nil
}
