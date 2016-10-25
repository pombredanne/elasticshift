package user

import (
	"strings"
	"time"

	"github.com/spf13/viper"
	"gitlab.com/conspico/esh/core/auth"
	"gitlab.com/conspico/esh/core/util"
	"gitlab.com/conspico/esh/team"
	"golang.org/x/crypto/bcrypt"
)

// Service ..
type Service interface {
	Create(teamName, domain, fullName, email, password string) (string, error)
	SignIn(teamName, domain, email, password string) (string, error)
	SignOut() (bool, error)
	Verify(code string) (bool, error)
}

type service struct {
	userDS Datastore
	teamDS team.Datastore
	config *viper.Viper
	signer interface{}
}

// NewService ..
func NewService(u Datastore, t team.Datastore, conf *viper.Viper, signer interface{}) Service {

	return &service{
		userDS: u,
		teamDS: t,
		config: conf,
		signer: signer,
	}
}

type verification struct {
	Team   string
	Email  string
	Expire time.Time
}

// Create a new user for a team
func (s service) Create(teamName, domain, fullname, email, password string) (string, error) {

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

	return s.generateAuthToken(teamID, email)
}

// SignIn ..
func (s service) SignIn(teamName, domain, email, password string) (string, error) {

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
	return s.generateAuthToken(teamID, user.ID)
}

// SignOut ..
func (s service) SignOut() (bool, error) {
	return true, nil
}

// Generates the auth token
func (s service) generateAuthToken(teamID, userID string) (string, error) {

	t := auth.Token{
		TeamID: teamID,
		UserID: userID,
	}

	signedStr, err := auth.GenerateToken(s.signer, t)
	if err != nil {
		return "", err
	}
	return signedStr, nil
}

// Verify .. the user code sent via email
func (s service) Verify(code string) (bool, error) {

	//teamID, err := s.teamRepository.GetTeamID(teamName)
	decrypted, err := util.XORDecrypt(s.config.GetString("key.verifycode"), code)
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

func (s service) getTeamID(teamName, domain string) (string, error) {

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
