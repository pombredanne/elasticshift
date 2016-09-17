package user

import (
	"crypto/rand"
	"fmt"
	"math/big"
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

	teamID, err := t.teamRepository.GetTeamID(teamName)

	result, err := t.userRepository.CheckExists(email, teamID)
	if result {
		return false, errUserAlreadyExists
	}

	userName := strings.Split(email, "@")[0]

	// generate verify code
	n, _ := rand.Int(rand.Reader, big.NewInt(999999))
	randCode := fmt.Sprintf("%06d", n.Int64())

	id, _ := util.NewUUID()

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
		VerifyCode: randCode,
		LastLogin:  time.Now(),
		CreatedBy:  "sysadmin",
		CreatedDt:  time.Now(),
		UpdatedBy:  "sysadmin",
		UpdatedDt:  time.Now(),
	}

	err = t.userRepository.Save(user)

	return true, err
}

// NewService ..
func NewService(u Repository, t team.Repository) Service {
	return &service{
		userRepository: u,
		teamRepository: t,
	}
}
