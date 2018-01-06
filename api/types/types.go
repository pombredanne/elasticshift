/*
Copyright 2017 The Elasticshift Authors.
*/
package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Client ..
type Client struct {
	ID           string   `bson:"_id,omitempty"`
	Secret       string   `bson:"secret"`
	Name         string   `bson:"name"`
	RedirectURIs []string `bson:"redirect_uris"`
	TrustedPeers []string `bson:"trusted_peers"`
	Public       bool     `bson:"public"`
	LogoURL      string   `bson:"logo_url"`
}

// GenericConf ..
type GenericSysConf struct {
	ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name  string        `json:"name" bson:"name"`
	Kind  string        `json:"kind" bson:"kind"`
	Value string        `json:"value" bson:"value"`
}

const (
	ReadOnly int = iota + 1
	ReadWrite
)

type NFSVolumeSysConf struct {
	ID         bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name       string        `json:"name" bson:"name"`
	Kind       string        `bson:"kind" bson:"kind,omitempty"`
	Server     string        `json:"server" bson:"server"`
	AccessMode int           `json:"access_mode" bson:"access_mode"`
}

// VCSSysConf ..(sysconf)
type VCSSysConf struct {
	// common fields for any sys config
	ID   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name string        `bson:"name,omitempty" json:"name"`
	Kind string        `bson:"kind,omitempty" json:"kind"`

	Key         string `bson:"key,omitempty" json:"key"`
	Secret      string `bson:"secret,omitempty" json:"secret"`
	CallbackURL string `bson:"callback_url,omitempty" json:"callback_url"`
	HookURL     string `bson:"hook_url,omitempty" json:"hook_url"`
}

// Team ..
type Team struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string        `bson:"name"`
	Display  string        `bson:"display,omitempty"`
	Accounts []VCS         `bson:"accounts"`
}

// User ..
type User struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Fullname      string        `bson:"fullname"`
	Username      string        `bson:"username"`
	Email         string        `bson:"email"`
	Password      string        `bson:"password"`
	Locked        bool          `bson:"locked"`
	Active        bool          `bson:"active"`
	BadAttempt    int8          `bson:"bad_attempt"`
	EmailVefified bool          `bson:"email_verified"`
	Scope         []string      `bson:"scope"`
	Team          string        `bson:"team"`
}

// VCS contains the information common amongst most OAuth and OAuth2 providers.
// All of the "raw" datafrom the provider can be found in the `RawData` field.
type VCS struct {
	ID           string    `json:"id" bson:"id,omitempty"`
	Name         string    `json:"name" bson:"name,omitempty"`
	Kind         string    `json:"kind" bson:"kind,omitempty"`
	Link         string    `json:"link" bson:"link,omitempty"`
	Source       string    `json:"source" bson:"source"`
	OwnerType    string    `json:"owner_type" bson:"owner_type,omitempty"`
	AvatarURL    string    `json:"avatar" bson:"avatar,omitempty"`
	AccessCode   string    `json:"access_code" bson:"access_code,omitempty"`
	AccessToken  string    `json:"access_token" bson:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token" bson:"refresh_token,omitempty"`
	TokenExpiry  time.Time `json:"token_expiry" bson:"token_expiry,omitempty"`
}

// Repository ..
// Represents as vcs repositories or projects
type Repository struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	RepoID        string        `json:"repo_id" bson:"repo_id"`
	VcsID         string        `json:"vcs_id" bson:"vcs_id"`
	Name          string        `json:"name" bson:"name,omitempty"`
	Private       bool          `json:"private" bson:"private,omitempty"`
	Link          string        `json:"link" bson:"link,omitempty"`
	Description   string        `json:"description" bson:"description,omitempty"`
	Fork          bool          `json:"fork" bson:"fork,omitempty"`
	DefaultBranch string        `json:"default_branch" bson:"default_branch,omitempty"`
	CloneURL      string        `json:"clone_url" bson:"clone_url"`
	Language      string        `json:"language" bson:"language,omitempty"`
	Identifier    string        `json:"-" bson:"identifier"`
	Team          string        `json:"-" bson:"team"`
}

type VCSList struct {
	Nodes []VCS `json:"nodes"`
	Count int   `json:"count"`
}

type RepositoryList struct {
	Nodes []Repository `json:"nodes"`
	Count int          `json:"count"`
}

//go:generate stringer -type=BuildStatus
type BuildStatus int

const (
	Stuck BuildStatus = iota + 1
	Running
	Success
	Failed
	Cancelled
	Waiting
)

func (b *BuildStatus) SetBSON(raw bson.Raw) error {

	var result int
	err := raw.Unmarshal(&result)
	if err != nil {
		return err
	}

	*b = BuildStatus(result)
	return nil
}

type Build struct {
	ID           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	RepositoryID string        `json:"repository_id" bson:"repository_id"`
	VcsID        string        `json:"vcs_id" bson:"vcs_id"`
	ContainerID  string        `json:"container_id" bson:"container_id"`
	Log          string        `json:"log" bson:"log"`
	StartedAt    time.Time     `json:"started_at" bson:"started_at"`
	EndedAt      time.Time     `json:"ended_at" bson:"ended_at"`
	TriggeredBy  string        `json:"triggered_by" bson:"triggered_by"`
	Status       BuildStatus   `json:"status" bson:"status"`
	Branch       string        `json:"branch" bson:"branch"`
	Team         string        `json:"team" bson:"team"`
}

type BuildList struct {
	Nodes []Build `json:"nodes"`
	Count int     `json:"count"`
}

type Container struct {
	ID              bson.ObjectId `json:"id" bson:"_id,omitempty"`
	BuildID         string        `json:"build_id" bson:"build_id"`
	RepoID          string        `json:"repo_id" bson:"repo_id"`
	VcsID           string        `json:"vcs_id" bson:"vcs_id"`
	OrchestrationID string        `json:"orchestration_id" bson:"orchestration_id"`
	Image           string        `json:"image" bson:"image"`
	repository_id   string        `json:"repository_id" bson:"repository_id"`
	StartedAt       time.Time     `json:"started_at" bson:"started_at"`
	StoppedAt       time.Time     `json:"stopped_at" bson:"stopped_at"`
	Duration        string        `json:"duration" bson:"duration"`
	Kind            string        `json:"kind" bson:"kind"`
}
