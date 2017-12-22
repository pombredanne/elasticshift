/*
Copyright 2017 The Elasticshift Authors.
*/
package team

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/core/utils"
)

var (
	// Team
	errTeamNameIsEmpty         = errors.New("Team name is empty")
	errTeamNameMinLength       = errors.New("Team name should should be minimum of 6 chars")
	errTeamNameMaxLength       = errors.New("Team name should not exceed 63 chars")
	errTeamNameContainsSymbols = errors.New("Team name should be alpha-numeric, no special chars or whitespace is allowed")
	errTeamAlreadyExists       = errors.New("Team name already exists")
	errTeamNameOrIdNeeded      = errors.New("Team name or ID must be given")
	errFailedToCreateTeam      = errors.New("Failed to create team")
)

type resolver struct {
	store  Store
	logger logrus.Logger
}

func (r *resolver) CreateTeam(params graphql.ResolveParams) (interface{}, error) {

	name := params.Args["name"].(string)

	// team name validation
	nameLength := len(name)
	if nameLength == 0 {
		return nil, errTeamNameIsEmpty
	}

	if !utils.IsAlphaNumericOnly(name) {
		return nil, errTeamNameContainsSymbols
	}

	if nameLength < 6 {
		return nil, errTeamNameMinLength
	}

	if nameLength > 63 {
		return nil, errTeamNameMaxLength
	}

	result, err := r.store.CheckExists(name)
	if result {
		return nil, errTeamAlreadyExists
	}

	t := &types.Team{
		Name:     name,
		Display:  name,
		Accounts: []types.VCS{},
	}

	err = r.store.Save(t)
	if err != nil {
		return nil, errFailedToCreateTeam
	}

	return t, err
}

func (r *resolver) FetchByNameOrID(params graphql.ResolveParams) (interface{}, error) {

	id, _ := params.Args["id"].(string)
	name, _ := params.Args["name"].(string)

	if id == "" && name == "" {
		return nil, errTeamNameOrIdNeeded
	}

	t, err := r.store.GetTeam(id, name)

	return &t, err
}

func (r *resolver) FetchTeams(params graphql.ResolveParams) (interface{}, error) {

	result := make([]types.Team, 1)

	return result, nil
}
