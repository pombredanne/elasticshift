/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/identity/team"
)

// VCSServer ..
//type VCSServer interface {
//authorize(w http.ResponseWriter, r *http.Request)
//handleAuthorizeCallback(w http.ResponseWriter, r *http.Request)
//GetVCS(teamID string) (GetVCSResponse, error)
//SyncVCS(r SyncVCSRequest) (bool, error)
//}

var (
	// VCS errors
	errFailedToFetchVCS = errors.New("Unknown vcs id")
)

type resolver struct {
	store     Store
	teamStore team.Store
	logger    logrus.Logger
}

func (r resolver) FetchVCSByTeamID(params graphql.ResolveParams) (interface{}, error) {

	teamName, _ := params.Args["team"].(string)
	if teamName == "" {
		return nil, team.ErrTeamNameIsEmpty
	}

	result, err := r.teamStore.GetVCS(teamName)

	var res types.VCSList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}

func (r resolver) FetchVCS(params graphql.ResolveParams) (interface{}, error) {

	result := make([]types.VCS, 1)

	return result, nil
}

// func (s resolver) Sync(ctx context.Context, req *api.SyncVCSReq) (*api.SyncVCSRes, error) {
// 	return nil, nil
// }
