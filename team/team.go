package team

import (
	"errors"
	"time"
)

var (
	errDomainNameIsEmpty         = errors.New("Team name is empty")
	errDomainNameMinLength       = errors.New("Team name should should be minimum of 6 chars")
	errDomainNameMaxLength       = errors.New("Team name should not exceed 63 chars")
	errDomainNameContainsSymbols = errors.New("Team name should be alpha-numeric, no special chars or whitespace is allowed")
	errTeamAlreadyExists         = errors.New("Team name already exists")
)

// Team ..
type Team struct {
	ID        string
	Domain    string
	Name      string
	CreatedDt time.Time
	CreatedBy string
	UpdatedDt time.Time
	UpdatedBy string
}

// Repository provides access a team.
type Repository interface {
	Save(team *Team) error
	CheckExists(name string) (bool, error)
	GetTeamID(name string) (string, error)
	FindByName(name string) (Team, error)
}
