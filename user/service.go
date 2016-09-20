package user

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gitlab.com/conspico/esh/core/util"
	"gitlab.com/conspico/esh/team"
)

// Service ..
type Service interface {
	Create(teamName, firstName, lastName, email string) (string, error)
	Verify(code string) (bool, error)
}

type service struct {
	userRepository Repository
	teamRepository team.Repository
	config         *viper.Viper
}

type verification struct {
	Team   string
	Email  string
	Code   string
	Expire time.Time
}

// Create a new user for a team
func (s service) Create(teamName, firstname, lastname, email string) (string, error) {

	// TODO : user teamid from session
	teamID, err := s.teamRepository.GetTeamID(teamName)

	result, err := s.userRepository.CheckExists(email, teamID)
	if result {
		return "", errUserAlreadyExists
	}

	// strip username from email
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

	err = s.userRepository.Save(user)

	v := verification{
		Team:   teamID,
		Email:  email,
		Expire: time.Now().AddDate(0, 0, 7), // a week
	}

	//send random code and link via email
	cipherText, err := util.EncryptStruct(s.config.GetString("key.verifycode"), v)
	if err != nil {
		return "", err
	}

	return cipherText, err
}

// Verify .. the user code sent via email
func (s service) Verify(code string) (bool, error) {

	//teamID, err := s.teamRepository.GetTeamID(teamName)
	var v verification
	err := util.DecryptStruct(s.config.GetString("key.verifycode"), code, v)
	if err != nil {
		return false, errVerificationCodeFailed
	}

	expireAt := v.Expire.Sub(time.Now())
	if expireAt.Hours() <= 0 && expireAt.Minutes() <= 0 {
		return false, errVerificationCodeExpired
	}

	return true, nil
}

// NewService ..
func NewService(u Repository, t team.Repository, conf *viper.Viper) Service {
	return &service{
		userRepository: u,
		teamRepository: t,
		config:         conf,
	}
}
