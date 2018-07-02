/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/pkg/identity/oauth2/providers"
	"gitlab.com/conspico/elasticshift/pkg/identity/team"
)

var (
	// VCS errors
	errFailedToFetchVCS = errors.New("Unknown vcs id")
	errNoURIProvided    = errors.New("URI is empty")
)

type resolver struct {
	store     Store
	teamStore team.Store
	logger    logrus.Logger
	providers providers.Providers
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
