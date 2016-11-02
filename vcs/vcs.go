package vcs

import "time"

var (
	errNoProviderFound       = "No provider found for %s"
	errGetUpdatedFokenFailed = "Failed to get updated token %s"
	errGettingRepositories   = "Failed to get repositories for %s"
)

// VCS user type
const (
	GithubType    = 1
	GitlabType    = 2
	BitBucketType = 3
	SvnType       = 4
	TfsType       = 5
)

// VCS account owner type
const (
	OwnerTypeUser = 1
	OwnerTypeOrg  = 2
)

// Datastore provides access a user.
type Datastore interface {

	// VCS account related operations
	Save(user *VCS) error
	GetVCS(teamID string) ([]VCS, error)
	GetByID(id string) (VCS, error)
	UpdateVCS(old *VCS, updated VCS) error
}

// VCS contains the information common amongst most OAuth and OAuth2 providers.
// All of the "raw" datafrom the provider can be found in the `RawData` field.
type VCS struct {
	ID           string `json:"ID"`
	TeamID       string `json:"-"`
	Name         string
	Type         int
	OwnerType    int
	AvatarURL    string    `json:"avatar"`
	AccessCode   string    `json:"-"`
	AccessToken  string    `json:"-"`
	RefreshToken string    `json:"-"`
	TokenType    string    `json:"-"`
	TokenExpiry  time.Time `json:"-"`
	CreatedDt    time.Time `json:"-"`
	CreatedBy    string    `json:"-"`
	UpdatedDt    time.Time `json:"lastUpdated"`
	UpdatedBy    string    `json:"-"`
}

// Repo ..
// Represents as vcs repositories or projects
type Repo struct {
	ID            string `json:"ID"`
	TeamID        string `json:"-"`
	VcsID         string `json:"-"`
	RepoID        string
	Name          string
	Private       string
	Link          string
	Description   string
	Fork          string
	DefaultBranch string
	Language      string
	CreatedDt     time.Time `json:"-"`
	CreatedBy     string    `json:"-"`
	UpdatedDt     time.Time `json:"lastUpdated"`
	UpdatedBy     string    `json:"-"`
}
