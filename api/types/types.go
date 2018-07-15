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
	SecretID     string    `json:"-" bson:"secret_id"`
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
	Source        string        `json:"source" bson:"source"`
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
	ID                bson.ObjectId `json:"id" bson:"_id,omitempty"`
	RepositoryID      string        `json:"repository_id" bson:"repository_id"`
	VcsID             string        `json:"vcs_id" bson:"vcs_id"`
	ContainerID       string        `json:"container_id" bson:"container_id"`
	Log               []Log         `json:"-" bson:"log"`
	LogType           string        `json:"-" bson:"log_type"`
	StartedAt         time.Time     `json:"started_at" bson:"started_at"`
	EndedAt           time.Time     `json:"ended_at" bson:"ended_at"`
	TriggeredBy       string        `json:"triggered_by" bson:"triggered_by"`
	Status            BuildStatus   `json:"status" bson:"status"`
	Branch            string        `json:"branch" bson:"branch"`
	CloneURL          string        `json:"clone_url" bson:"clone_url"`
	Language          string        `json:"language" bson:"language"`
	Team              string        `json:"team" bson:"team"`
	ContainerEngineID string        `json:"-" bson:"container_engine_id"`
	StorageID         string        `json:"-" bson:"storage_id"`
	StoragePath       string        `json:"-" bson:"storage_path"`
	Privatekey        string        `json:"-" bson:"private_key,omitempty"`
	Graph             string        `json:"graph" bson:"graph,omitempty"`
	Source            string        `json:"source" bson:"source"`
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

type Plugin struct {
	ID             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name           string        `json:"name" bson:"name"`
	Description    string        `json:"description" bson:"description"`
	Language       string        `json:"language" bson:"language"`
	Version        string        `json:"version" bson:"version"`
	Author         string        `json:"author" bson:"author"`
	Email          string        `json:"email" bson:"email"`
	SourceURL      string        `json:"source_url" bson:"source_url,omitempty"`
	Readme         string        `json:"readme" bson:"readme"`
	UsedTeamCount  int64         `json:"used_team_count" bson:"used_team_count"`
	UsedReposCount int64         `json:"used_build_count" bson:"used_build_count"`
	IconURL        string        `json:"icon_url" bson:"icon_url"`
	Ratings        string        `json:"ratings" bson:"ratings"`
	Team           string        `json:"team" bson:"team"`
	Path           string        `json:"path" bson:"path"`
}

type PluginList struct {
	Nodes []Plugin `json:"nodes"`
	Count int      `json:"count"`
}

type ContainerEngine struct {
	ID           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name         string        `json:"name" bson:"name"`
	Provider     int           `json:"provider" bson:"provider"`
	Kind         int           `json:"kind" bson:"kind"`
	Host         string        `json:"host" bson:"host,omitempty"`
	Certificate  string        `json:"certificate" bson:"certificate,omitempty"`
	Token        string        `json:"token" bson:"token,omitempty"`
	InternalType int           `json:"-" bson:"internal_type"`
	Team         string        `json:"team" bson:"team"`
	KubeFile     KubeConfig    `json:"kube_config" bson:"kube_config,omitempty"`
}

type ContainerEngineList struct {
	Nodes []ContainerEngine `json:"nodes"`
	Count int               `json:"count"`
}

type MinioStorage struct {
	Host        string `json:"host" bson:"host"`
	Certificate string `json:"certificate" bson:"certificate"`
	AccessKey   string `json:"access_key" bson:"access_key"`
	SecretKey   string `json:"secret_key" bson:"secret_key"`
}

type NFSStorage struct {
	Server    string `json:"server" bson:"server"`
	Path      string `json:"path" bson:"path"`
	ReadOnly  bool   `json:"readonly" bson:"read_only"`
	MountPath string `json:"mount_path" bson:"mount_path"`
}

type StorageSource struct {
	NFS   *NFSStorage   `json:"nfs" bson:"nfs,omitempty"`
	Minio *MinioStorage `json:"minio" bson:"minio,omitempty"`
}

type Storage struct {
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name          string        `json:"name" bson:"name"`
	Provider      int           `json:"provider" bson:"provider"`
	Kind          int           `json:"kind" bson:"kind"`
	InternalType  int           `json:"-" bson:"internal_type"`
	Team          string        `json:"team" bson:"team"`
	StorageSource `json:"storage_source" bson:"storage_source"`
}

type StorageList struct {
	Nodes []Storage `json:"nodes"`
	Count int       `json:"count"`
}

type Infrastructure struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description,omitempty"`
	Kind        string        `json:"kind" bson:"kind"`
	Private     bool          `json:"private" bson:"private"`
	Code        string        `json:"code" bson:"code"`
	Team        string        `json:"team" bson:"team"`
}

type InfrastructureList struct {
	Nodes []Infrastructure `json:"nodes"`
	Count int              `json:"count"`
}

type Secret struct {
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name          string        `json:"name" bson:"name"`
	Kind          string        `json:"kind" bson:"kind"`
	ReferenceKind string        `json:"reference_kind" bson:"reference_kind"`
	ReferenceID   string        `json:"reference_id" bson:"reference_id"`
	Value         string        `json:"value" bson:"value"`
	InternalType  string        `json:"-" bson:"internal_type"`
	KeyID         string        `json:"-" bson:"key_id"`
	TeamID        string        `json:"team_id" bson:"team_id"`
}

type SecretList struct {
	Nodes []Secret `json:"nodes"`
	Count int      `json:"count"`
}

type Property struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

type Default struct {
	ID                bson.ObjectId     `json:"id" bson:"_id,omitempty"`
	Kind              int               `json:"kind" bson:"kind"`
	ReferenceID       string            `json:"reference_id" bson:"reference_id"`
	ContainerEngineID string            `json:"container_engine_id" bson:"container_engine_id,omitempty"`
	StorageID         string            `json:"storage_id" bson:"storage_id,omitempty"`
	Languages         map[string]string `json:"languages" bson:"languages,omitempty"`
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

type Shiftfile struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	File        ShiftfileType `json:"file" bson:"file"`
	UsedByTeams int64         `json:"used_by_teams" bson:"used_by_teams"`
	UsedByRepos int64         `json:"used_by_repos" bson:"used_by_repos"`
	TeamID      string        `json:"team_id" bson:"team_id"`
	Ratings     string        `json:"ratings" bson:"ratings"`
}

type ShiftfileType []byte

func (f ShiftfileType) GetBSON() (interface{}, error) {
	return base64.StdEncoding.EncodeToString([]byte(f)), nil
}

func (f *ShiftfileType) SetBSON(raw bson.Raw) error {

	var data string
	var err error

	if err = raw.Unmarshal(&data); err != nil {
		return err
	}

	*f, err = base64.StdEncoding.DecodeString(data)
	return err
}

type ShiftfileList struct {
	Nodes []Shiftfile `json:"nodes"`
	Count int         `json:"count"`
}
