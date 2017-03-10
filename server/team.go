// Package server ..
// Teamor Ghazni Nattarshah
// Date: 2/11/17
package server

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/backup/core/util"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/store"
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

type teamServer struct {
	store  store.TeamStore
	logger logrus.FieldLogger
}

// NewTeamServer ..
// Implementation of api.UserServer
func NewTeamServer(s *Server) api.TeamServer {
	return &teamServer{
		store:  store.NewTeamStore(s.Store),
		logger: s.Logger,
	}
}

// Teamorize ..
func (s *teamServer) Create(ctx context.Context, req *api.CreateTeamReq) (*api.CreateTeamRes, error) {

	res := &api.CreateTeamRes{}

	// team name validation
	nameLength := len(req.Name)
	s.logger.Warnln(req.Name)
	if nameLength == 0 {
		return res, errTeamNameIsEmpty
	}

	if !util.IsAlphaNumericOnly(req.Name) {
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

	team := &store.Team{
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
