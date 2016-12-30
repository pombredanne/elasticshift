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
func (s userService) Create(teamName, domain, fullname, email, password string) (string, error) {

	tname := teamName
	if tname == "" {
		tname = domain
	}

	result, err := s.userDS.CheckExists(email, tname)
	if result {
		return "", errUserAlreadyExists
	}

	// strip username from email
	userName := strings.Split(email, "@")[0]

	// generate hashed password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", stacktrace.Propagate(err, errUserCreationFailed.Error())
	}

	user := &User{
		Fullname:      fullname,
		Username:      userName,
		Email:         email,
		Password:      string(hashedPwd[:]),
		Locked:        Unlocked,
		Active:        Active,
		BadAttempt:    0,
		EmailVefified: false,
		Team:          tname,
		Scope:         []string{},
	}

	err = s.userDS.Save(user)
	if err != nil {
		return "", stacktrace.Propagate(err, errUserCreationFailed.Error())
	}

	return s.generateAuthToken(tname, email, userName, tname)
}

// SignIn ..
func (s userService) SignIn(teamName, domain, email, password string) (string, error) {

	tname := teamName
	if tname == "" {
		tname = domain
	}

	user, err := s.userDS.GetUser(email, tname)
	if err != nil {
		return "", errInvalidEmailOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errInvalidEmailOrPassword.Error(), nil
	}
	return s.generateAuthToken(tname, user.ID.String(), user.Username, tname)
}

// SignOut ..
func (s userService) SignOut() (bool, error) {
	return true, nil
}

// Generates the auth token
func (s userService) generateAuthToken(teamID, userID, userName, teamName string) (string, error) {

	t := auth.Token{
		TeamID:   teamID,
		UserID:   userID,
		Username: userName,
		Teamname: teamName,
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
