package user

import (
	"strings"
	"time"

	"gitlab.com/conspico/esh/core/util"
	"gitlab.com/conspico/esh/team"
)

// Service ..
type Service interface {
	Create(teamName, firstName, lastName, email string) (bool, error)
}

type service struct {
	userRepository Repository
	teamRepository team.Repository
}

func (t service) Create(teamName, firstname, lastname, email string) (bool, error) {

	result, err := t.userRepository.CheckExists(email)
	if result {
		return false, errUserAlreadyExists
	}

	id, err := util.NewUUID()
	if err != nil {

	}

	teamID, err := t.teamRepository.GetTeamID(teamName)

	userName := strings.Split(email, "@")[0]

	user := &User{
		PUUID:      id,
		TeamID:     teamID,
		FirstName:  firstname,
		LastName:   lastname,
		UserName:   userName,
		Email:      email,
		Locked:     1,
		Active:     1,
		BadAttemt:  0,
		VerifyCode: 123456,
		LastLogin:  time.Now(),
		CreatedBy:  "sysadmin",
		CreatedDt:  time.Now(),
		UpdatedBy:  "sysadmin",
		UpdatedDt:  time.Now(),
	}

	t.userRepository.Save(user)

	return true, nil
}

// NewService ..
func NewService(u Repository, t team.Repository) Service {
	return &service{
		userRepository: u,
		teamRepository: t,
	}
}
