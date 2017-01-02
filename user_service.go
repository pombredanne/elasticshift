// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

import (
	"strings"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/palantir/stacktrace"
	"gitlab.com/conspico/esh/core/auth"
	"gitlab.com/conspico/esh/core/util"
	"golang.org/x/crypto/bcrypt"
)

const (

	// Separator used to delimit user info sent for activation
	Separator = ";"

	// Inactive ..
	Inactive = 0

	// Active ..
	Active = 1

	// Unlocked ..
	Unlocked = 0

	// Locked ..
	Locked = 1
)

type userService struct {
	userDS UserDatastore
	teamDS TeamDatastore
	config Config
	logger *logrus.Logger
	signer interface{}
}

// NewUserService ..
func NewUserService(ctx AppContext) UserService {

	return &userService{
		userDS: ctx.UserDatastore,
		teamDS: ctx.TeamDatastore,
		config: ctx.Config,
		signer: ctx.Signer,
		logger: ctx.Logger,
	}
}

type verification struct {
	Team   string
	Email  string
	Expire time.Time
}

// Create a new user for a team
func (s userService) Create(r signupRequest) (string, error) {

	//teamName, domain, fullname, email, password string
	result, err := s.userDS.CheckExists(r.Email, r.Team)
	if result {
		return "", errUserAlreadyExists
	}

	// strip username from email
	userName := strings.Split(r.Email, "@")[0]

	// generate hashed password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", stacktrace.Propagate(err, errUserCreationFailed.Error())
	}

	user := &User{
		Fullname:      r.Fullname,
		Username:      userName,
		Email:         r.Email,
		Password:      string(hashedPwd[:]),
		Locked:        Unlocked,
		Active:        Active,
		BadAttempt:    0,
		EmailVefified: false,
		Team:          r.Team,
		Scope:         []string{},
	}

	err = s.userDS.Save(user)
	if err != nil {
		return "", stacktrace.Propagate(err, errUserCreationFailed.Error())
	}

	return s.generateAuthToken(r.Email, userName, r.Team)
}

// SignIn ..
func (s userService) SignIn(r signInRequest) (string, error) {


	user, err := s.userDS.GetUser(r.Email, r.Team)
	if err != nil {
		return "", errInvalidEmailOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password))
	if err != nil {
		return errInvalidEmailOrPassword.Error(), nil
	}
	return s.generateAuthToken(user.ID.String(), user.Username, r.Team)
}

// SignOut ..
func (s userService) SignOut() (bool, error) {
	return true, nil
}

// Generates the auth token
func (s userService) generateAuthToken(userID, userName, team string) (string, error) {

	t := auth.Token{
		UserID:   userID,
		Username: userName,
		Team: team,
	}

	signedStr, err := auth.GenerateToken(s.signer, t)
	if err != nil {
		return "", stacktrace.Propagate(err, "Failed to geenrate auth token")
	}
	return signedStr, nil
}

// Verify .. the user code sent via email
func (s userService) Verify(code string) (bool, error) {

	//teamID, err := s.teamRepository.GetTeamID(teamName)
	decrypted, err := util.XORDecrypt(s.config.Key.VerifyCode, code)
	if err != nil {
		return false, errVerificationCodeFailed
	}

	// TODO fetch based on name and email and see if the data's been tampered
	v := strings.Split(decrypted, Separator)
	expireAt, err := time.Parse(time.RFC3339Nano, v[2])
	diff := expireAt.Sub(time.Now())

	if diff.Hours() <= 0 && diff.Minutes() <= 0 {
		return false, errVerificationCodeExpired
	}

	return true, nil
}
