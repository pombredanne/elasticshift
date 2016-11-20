package esh

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/palantir/stacktrace"
	"gitlab.com/conspico/esh/core/util"
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

	id, _ := util.NewUUID()
	if err != nil {
		return false, stacktrace.Propagate(err, "Unable to create team"+name)
	}

	team := &Team{
		ID:        id,
		Name:      name,
		Domain:    name,
		CreatedBy: "sysadmin",
		CreatedDt: time.Now(),
		UpdatedBy: "sysadmin",
		UpdatedDt: time.Now(),
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
