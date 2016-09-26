package user

import (
	"crypto/rand"
	"fmt"
	"math/big"
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
	Create(teamName, firstName, lastName, email, password string) (string, error)
	SignIn(team, email, password string) (string, error)
	Verify(code string) (bool, error)
}

type service struct {
	userRepository Repository
	teamRepository team.Repository
	config         *viper.Viper
	signer         []byte
}

// NewService ..
func NewService(u Repository, t team.Repository, conf *viper.Viper, signer []byte) Service {

	return &service{
		userRepository: u,
		teamRepository: t,
		config:         conf,
		signer:         signer,
	}
}

type verification struct {
	Team   string
	Email  string
	Expire time.Time
}

// Create a new user for a team
func (s service) Create(teamName, firstname, lastname, email, password string) (string, error) {

	// TODO : user teamid from session
	teamID, err := s.teamRepository.GetTeamID(teamName)
	if err != nil {
		return "", errNoTeamIDNotExist
	}

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

	// generate hashed password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errUserCreationFailed
	}

	user := &User{
		PUUID:        id,
		TeamID:       teamID,
		FirstName:    firstname,
		LastName:     lastname,
		UserName:     userName,
		Email:        email,
		Locked:       Unlocked,
		Active:       Active,
		BadAttempt:   0,
		PasswordHash: string(hashedPwd[:]),
		VerifyCode:   randCode,
		LastLogin:    time.Now(),
		CreatedBy:    "sysadmin",
		CreatedDt:    time.Now(),
		UpdatedBy:    "sysadmin",
		UpdatedDt:    time.Now(),
	}

	err = s.userRepository.Save(user)
	if err != nil {
		return "", errUserCreationFailed
	}

	return s.generateAuthToken(teamID, email)
}

func (s service) generateAuthToken(teamID, emailID string) (string, error) {

	t := auth.Token{
		Email:  emailID,
		TeamID: teamID,
	}

	signedStr, err := auth.GenerateToken(s.signer, t)
	if err != nil {
		return "", err
	}
	return signedStr, nil
}

// SignIn ..
func (s service) SignIn(team, email, password string) (string, error) {

	teamID, err := s.teamRepository.GetTeamID(team)
	if err != nil {
		return "", errNoTeamIDNotExist
	}
	user, err := s.userRepository.GetUser(email, teamID)
	if err != nil {
		return "", errInvalidEmailOrPassword
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errInvalidEmailOrPassword
	}
	return s.generateAuthToken(teamID, email)
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
	teamID := v[0]
	email := v[1]
	expireAt, err := time.Parse(time.RFC3339Nano, v[2])
	diff := expireAt.Sub(time.Now())
	fmt.Println("Team id = ", teamID)
	fmt.Println("Email id = ", email)
	fmt.Println(diff)

	if diff.Hours() <= 0 && diff.Minutes() <= 0 {
		return false, errVerificationCodeExpired
	}

	return true, nil
}
