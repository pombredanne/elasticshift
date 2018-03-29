/*
Copyright 2017 The Elasticshift Authors.
*/
package types

import (
	"encoding/base64"
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
	Path       string        `json:"path" bson:"path"`
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
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string        `bson:"name"`
	Display    string        `bson:"display,omitempty"`
	Accounts   []VCS         `bson:"accounts"`
	KubeConfig KubeConfig    `json:"-" bson:"kube_config"`
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
	BS_STUCK BuildStatus = iota + 1
	BS_RUNNING
	BS_SUCCESS
	BS_FAILED
	BS_CANCELLED
	BS_WAITING
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

type Log struct {
	Time time.Time `json:"time" bson:"time"`
	Data string    `json:"data" bson:"data"`
}

type Build struct {
	ID           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	RepositoryID string        `json:"repository_id" bson:"repository_id"`
	VcsID        string        `json:"vcs_id" bson:"vcs_id"`
	ContainerID  string        `json:"container_id" bson:"container_id"`
	Log          []Log         `json:"-" bson:"log"`
	LogType      string        `json:"-" bson:"log_type"`
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

//go:generate stringer -type=ContainerStatus
type ContainerStatus int

const (
	CS_RUNNING ContainerStatus = iota + 1
	CS_STARTED
	CS_STOPPED
)

func (b *ContainerStatus) SetBSON(raw bson.Raw) error {

	var result int
	err := raw.Unmarshal(&result)
	if err != nil {
		return err
	}

	*b = ContainerStatus(result)
	return nil
}

type Container struct {
	ID              bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	BuildID         string          `json:"build_id" bson:"build_id"`
	RepositoryID    string          `json:"repository_id" bson:"repository_id"`
	ContainerID     string          `json:"container_id" bson:"container_id"`
	VcsID           string          `json:"vcs_id" bson:"vcs_id"`
	OrchestrationID string          `json:"orchestration_id" bson:"orchestration_id"`
	Image           string          `json:"image" bson:"image"`
	StartedAt       time.Time       `json:"started_at" bson:"started_at"`
	StoppedAt       time.Time       `json:"stopped_at" bson:"stopped_at"`
	Duration        string          `json:"duration" bson:"duration"`
	Kind            string          `json:"kind" bson:"kind"`
	Status          ContainerStatus `json:"status" bson:"status"`
}

type ContainerList struct {
	Nodes []Container `json:"nodes"`
	Count int         `json:"count"`
}

type App struct {
	ID             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name           string        `json:"name" bson:"name"`
	Description    string        `json:"description" bson:"description"`
	Language       string        `json:"language" bson:"language"`
	Version        string        `json:"version" bson:"version"`
	UsedTeamCount  int64         `json:"used_team_count" bson:"used_team_count"`
	UsedBuildCount int64         `json:"used_build_count" bson:"used_build_count"`
	IconURL        string        `json:"icon_url" bson:"icon_url"`
	SourceURL      string        `json:"source_url" bson:"source_url"`
	Readme         string        `json:"readme" bson:"readme"`
	Ratings        string        `json:"ratings" bson:"ratings"`
}

type AppList struct {
	Nodes []App `json:"nodes"`
	Count int   `json:"count"`
}

type KubeConfig []byte

func (f KubeConfig) GetBSON() (interface{}, error) {
	return base64.StdEncoding.EncodeToString([]byte(f)), nil
}

func (f *KubeConfig) SetBSON(raw bson.Raw) error {

	var data string
	var err error

	if err = raw.Unmarshal(&data); err != nil {
		return err
	}

	*f, err = base64.StdEncoding.DecodeString(data)
	return err
}
