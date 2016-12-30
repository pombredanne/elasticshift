// Package esh ...
// Author: Ghazni Nattarshah
// Date: Oct 29, 2016
package esh

import (
	"github.com/Sirupsen/logrus"
	"github.com/palantir/stacktrace"
)

type teamservice struct {
	teamDS TeamDatastore
	logger *logrus.Logger
}

func (s teamservice) Create(name string) (bool, error) {

	result, err := s.teamDS.CheckExists(name)
	if result {
		return false, errTeamAlreadyExists
	}

	team := &Team{
		Name:    name,
		Display: name,
	}

	err = s.teamDS.Save(team)

	return err == nil, stacktrace.Propagate(err, "Unable to create team")
}

// NewTeamService ..
func NewTeamService(ctx AppContext) TeamService {
	return &teamservice{
		teamDS: ctx.TeamDatastore,
		logger: ctx.Logger,
	}
}
