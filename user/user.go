package user

import (
	"errors"
	"time"
)

var (
	errUserAlreadyExists         = errors.New("User already exists")
	errUserCreationFailed        = errors.New("User creation failed ")
	errNoTeamIDNotExist          = errors.New("Team ID does not exist")
	errNOTeamIDProvided          = errors.New("No team id provided")
	errVerificationCodeIsEmpty   = errors.New("Verification code is empty")
	errVerificationCodeExpired   = errors.New("Verification code expired")
	errVerificationCodeFailed    = errors.New("Verification code seems to be failed")
	errUsernameOrPasswordIsEmpty = errors.New("User or password is empty ")
	errInvalidEmailOrPassword    = errors.New("Invalid email or password")
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

// User ..
type User struct {
	ID         string
	TeamID     string
	Fullname   string
	Username   string
	Email      string
	Password   string
	Locked     int8
	Active     int8
	BadAttempt int8
	CreatedDt  time.Time
	CreatedBy  string
	UpdatedDt  time.Time
	UpdatedBy  string
}

// Repository provides access a user.
type Repository interface {
	Save(user *User) error
	CheckExists(email, teamID string) (bool, error)
	GetUser(email, teamID string) (User, error)
	FindByName(username string) (User, error)
}
