/*
Copyright 2017 The Elasticshift Authors.
*/
package team

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/api/types"
	core "gitlab.com/conspico/elasticshift/core/store"
	"gitlab.com/conspico/elasticshift/core/utils"
)

var (
	// Team
	errTeamNameIsEmpty         = errors.New("Team name is empty")
	errTeamNameMinLength       = errors.New("Team name should should be minimum of 6 chars")
	errTeamNameMaxLength       = errors.New("Team name should not exceed 63 chars")
	errTeamNameContainsSymbols = errors.New("Team name should be alpha-numeric, no special chars or whitespace is allowed")
	errTeamAlreadyExists       = errors.New("Team name already exists")
	errFailedToCreateTeam      = errors.New("Failed to create team")
)

type server struct {
	store  Store
	logger logrus.FieldLogger
}

// NewServer ..
// Implementation of api.UserServer
func NewServer(s core.Store, logger logrus.FieldLogger) api.TeamServer {
	return &server{
		store:  NewStore(s),
		logger: logger,
	}
}

// Teamorize ..
func (s *server) Create(ctx context.Context, req *api.CreateTeamReq) (*api.CreateTeamRes, error) {

	res := &api.CreateTeamRes{}

	// team name validation
	nameLength := len(req.Name)
	if nameLength == 0 {
		return res, errTeamNameIsEmpty
	}

	if !utils.IsAlphaNumericOnly(req.Name) {
		return res, errTeamNameContainsSymbols
	}

	if nameLength < 6 {
		return res, errTeamNameMinLength
	}

	if nameLength > 63 {
		return res, errTeamNameMaxLength
	}

	result, err := s.store.CheckExists(req.Name)
	if result {

		res.AlreadyExist = true
		return res, errTeamAlreadyExists
	}

	team := &types.Team{
		Name:    req.Name,
		Display: req.Name,
	}

	err = s.store.Save(team)
	if err != nil {
		return res, errFailedToCreateTeam
	}

	res.Created = true
	return res, err
}
