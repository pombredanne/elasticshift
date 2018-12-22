/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
	"context"
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	"github.com/elasticshift/elasticshift/internal/shiftserver/identity/oauth2/providers"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
	"github.com/elasticshift/elasticshift/internal/shiftserver/team"
)

var (
	// VCS errors
	errFailedToFetchVCS = errors.New("Unknown vcs id")
	errNoURIProvided    = errors.New("URI is empty")
)

// Resolver ...
type Resolver interface {
	FetchVCS(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	store     store.Vcs
	teamStore store.Team
	logger    *logrus.Entry
	providers providers.Providers
}

// NewResolver ...
func NewResolver(ctx context.Context, loggr logger.Loggr, s store.Shift, providers providers.Providers) (Resolver, error) {

	r := &resolver{
		store:     s.Vcs,
		teamStore: s.Team,
		logger:    loggr.GetLogger("graphql/vcs"),
		providers: providers,
	}
	return r, nil
}

func (r resolver) FetchVCS(params graphql.ResolveParams) (interface{}, error) {

	teamID, _ := params.Args["team"].(string)
	if teamID == "" {
		return nil, team.ErrTeamNameIsEmpty
	}

	result, err := r.teamStore.GetVCS(teamID)

	var res types.VCSList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}
