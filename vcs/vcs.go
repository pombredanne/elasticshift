package vcs

import "time"

var (
	errNoProviderFound = "No provider found for %s"
)

// VCS user type
const (
	GithubType    = 1
	GitlabType    = 2
	BitBucketType = 3
	SvnType       = 4
	TfsType       = 5
)

// VCS contains the information common amongst most OAuth and OAuth2 providers.
// All of the "raw" datafrom the provider can be found in the `RawData` field.
type VCS struct {
	ID          string `json:"ID"`
	TeamID      string `json:"-"`
	Name        string
	Type        int
	AccessCode  string    `json:"-"`
	AccessToken string    `json:"-"`
	AvatarURL   string    `json:"avatar"`
	CreatedDt   time.Time `json:"-"`
	CreatedBy   string    `json:"-"`
	UpdatedDt   time.Time `json:"lastUpdated"`
	UpdatedBy   string    `json:"-"`
}

// Repository provides access a user.
type Repository interface {
	Save(user *VCS) error
	GetVCS(teamID string) ([]VCS, error)
}
