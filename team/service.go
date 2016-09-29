package team

import (
	"time"

	"gitlab.com/conspico/esh/core/util"
)

// Service ..
type Service interface {
	Create(name string) (bool, error)
}

type service struct {
	teamRepository Repository
}

func (s service) Create(name string) (bool, error) {

	result, err := s.teamRepository.CheckExists(name)
	if result {
		return false, errTeamAlreadyExists
	}

	id, err := util.NewUUID()
	if err != nil {

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

	err = s.teamRepository.Save(team)

	return err == nil, err
}

// NewService ..
func NewService(t Repository) Service {
	return &service{
		teamRepository: t,
	}
}
