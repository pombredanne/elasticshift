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
	ID          string
	TeamID      string
	Name        string
	Type        int
	AccessToken string
	AvatarURL   string
	CreatedDt   time.Time
	CreatedBy   string
	UpdatedDt   time.Time
	UpdatedBy   string
}

/*ID VARCHAR(32) PRIMARY KEY COMMENT 'UUID of vcs account',
  TEAM_ID VARCHAR(32) COMMENT 'UUID of TEAM',
  NAME VARCHAR(100) NOT NULL COMMENT 'username of version control system',
  TYPE TINYINT DEFAULT 0 COMMENT '1-github, 2-gitlab, 3-bitbucket, 4-SVN, 5-TFS',
  ACCESS_CODE VARCHAR(255) NOT NULL COMMENT 'Hashed access code',
  AVATAR_URL VARCHAR(255) NOT NULL COMMENT 'user avatar',
  REVOKED TINYINT DEFAULT 0 COMMENT '1-revoked, 0-active',
  CREATED_DT DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT 'VCS account creation datetime',
  UPDATED_DT DATETIME ON UPDATE CURRENT_TIMESTAMP COMMENT 'VCS account updated datetime',
  CREATED_BY VARCHAR(100) COMMENT 'who created this record',
  UPDATED_BY */

// Repository provides access a user.
type Repository interface {
	Save(user *VCS) error
}
