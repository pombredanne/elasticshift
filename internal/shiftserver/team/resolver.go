/*
Copyright 2017 The Elasticshift Authors.
*/
package team

import (
	"context"
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
)

var (
	// Team
	ErrTeamNameIsEmpty         = errors.New("Team name is empty")
	ErrTeamNameMinLength       = errors.New("Team name should should be minimum of 6 chars")
	ErrTeamNameMaxLength       = errors.New("Team name should not exceed 63 chars")
	ErrTeamNameContainsSymbols = errors.New("Team name should be alpha-numeric, no special chars or whitespace is allowed")
	ErrTeamAlreadyExists       = errors.New("Team name already exists")
	ErrTeamNameOrIdNeeded      = errors.New("Team name or ID must be given")
	ErrFailedToCreateTeam      = errors.New("Failed to create team")
)

// Resolver ...
type Resolver interface {
	CreateTeam(params graphql.ResolveParams) (interface{}, error)
	FetchByNameOrID(params graphql.ResolveParams) (interface{}, error)
	FetchTeams(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	store  store.Team
	logger *logrus.Entry
}

// NewResolver ...
func NewResolver(ctx context.Context, loggr logger.Loggr, s store.Shift) (Resolver, error) {

	r := &resolver{
		store:  s.Team,
		logger: loggr.GetLogger("graphql/team"),
	}
	return r, nil
}

func (r *resolver) CreateTeam(params graphql.ResolveParams) (interface{}, error) {

	name := params.Args["name"].(string)

	// team name validation
	nameLength := len(name)
	if nameLength == 0 {
		return nil, ErrTeamNameIsEmpty
	}

	if !utils.IsAlphaNumericOnly(name) {
		return nil, ErrTeamNameContainsSymbols
	}

	if nameLength < 6 {
		return nil, ErrTeamNameMinLength
	}

	if nameLength > 63 {
		return nil, ErrTeamNameMaxLength
	}

	result, err := r.store.CheckExists(name)
	if result {
		return nil, ErrTeamAlreadyExists
	}

	t := &types.Team{
		Name:     name,
		Display:  name,
		Accounts: []types.VCS{},
	}

	err = r.store.Save(t)
	if err != nil {
		return nil, ErrFailedToCreateTeam
	}

	return t, err
}

func (r *resolver) FetchByNameOrID(params graphql.ResolveParams) (interface{}, error) {

	id, _ := params.Args["id"].(string)
	name, _ := params.Args["name"].(string)

	if id == "" && name == "" {
		return nil, ErrTeamNameOrIdNeeded
	}

	t, err := r.store.GetTeam(id, name)

	return &t, err
}

func (r *resolver) FetchTeams(params graphql.ResolveParams) (interface{}, error) {

	result := make([]types.Team, 1)

	return result, nil
}
