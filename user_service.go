package esh

import (
	"strings"
	"time"

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
	signer interface{}
}

// NewUserService ..
func NewUserService(appCtx AppContext) UserService {

	return &userService{
		userDS: appCtx.UserDatastore,
		teamDS: appCtx.TeamDatastore,
		config: appCtx.Config,
		signer: appCtx.Signer,
	}
}

type verification struct {
	Team   string
	Email  string
	Expire time.Time
}

// Create a new user for a team
func (s userService) Create(teamName, domain, fullname, email, password string) (string, error) {

	teamID, err := s.getTeamID(teamName, domain)
	if err != nil {
		return "", err
	}

	result, err := s.userDS.CheckExists(email, teamID)
	if result {
		return "", errUserAlreadyExists
	}

	// strip username from email
	userName := strings.Split(email, "@")[0]

	id, _ := util.NewUUID()

	// generate hashed password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errUserCreationFailed
	}

	user := &User{
		ID:         id,
		TeamID:     teamID,
		Fullname:   fullname,
		Username:   userName,
		Email:      email,
		Locked:     Unlocked,
		Active:     Active,
		BadAttempt: 0,
		Password:   string(hashedPwd[:]),
		CreatedBy:  "sysadmin",
		CreatedDt:  time.Now(),
		UpdatedBy:  "sysadmin",
		UpdatedDt:  time.Now(),
	}

	err = s.userDS.Save(user)
	if err != nil {
		return "", errUserCreationFailed
	}

	tname := teamName
	if tname == "" {
		tname = domain
	}
	return s.generateAuthToken(teamID, email, userName, tname)
}

// SignIn ..
func (s userService) SignIn(teamName, domain, email, password string) (string, error) {

	teamID, err := s.getTeamID(teamName, domain)
	if err != nil {
		return "", err
	}
	user, err := s.userDS.GetUser(email, teamID)
	if err != nil {
		return "", errInvalidEmailOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errInvalidEmailOrPassword.Error(), nil
	}
	tname := teamName
	if tname == "" {
		tname = domain
	}
	return s.generateAuthToken(teamID, user.ID, user.Username, tname)
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
		return "", err
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

func (s userService) getTeamID(teamName, domain string) (string, error) {

	// checks the team from subdomain
	teamID, err := s.teamDS.GetTeamID(domain)
	if err != nil {
		// checks the team from JSON request
		teamID, err = s.teamDS.GetTeamID(teamName)
		if err != nil {
			return "", errNoTeamIDNotExist
		}
	}
	return teamID, nil
}
