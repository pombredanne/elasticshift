package user

import (
	"errors"
	"time"
)

var (
	errUserAlreadyExists       = errors.New("User already exists")
	errNoTeamIDNotExist        = errors.New("Team ID does not exist")
	errVerificationCodeIsEmpty = errors.New("Verification code is empty")
	errVerificationCodeExpired = errors.New("Verification code expired")
	errVerificationCodeFailed  = errors.New("Verification code seems to be failed")
)

// User ..
type User struct {
	PUUID        string `gorm:"column:puuid"`
	TeamID       string `gorm:"column:team_puuid"`
	FirstName    string `gorm:"column:firstname"`
	LastName     string `gorm:"column:lastname"`
	UserName     string `gorm:"column:username"`
	Email        string
	PasswordHash string `gorm:"column:hashed_password"`
	Locked       int8
	Active       int8
	BadAttemt    int8      `gorm:"column:bad_attempt"`
	LastLogin    time.Time `gorm:"column:last_login"`
	VerifyCode   string    `gorm:"column:verify_code"`
	CreatedDt    time.Time
	CreatedBy    string
	UpdatedDt    time.Time
	UpdatedBy    string
}

// Repository provides access a user.
type Repository interface {
	Save(user *User) error
	CheckExists(email, teamID string) (bool, error)
	FindByName(username string) (User, error)
}
