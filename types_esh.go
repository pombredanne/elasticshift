// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

import (
	"context"
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

var (
	// Team
	errDomainNameIsEmpty         = errors.New("Team name is empty")
	errDomainNameMinLength       = errors.New("Team name should should be minimum of 6 chars")
	errDomainNameMaxLength       = errors.New("Team name should not exceed 63 chars")
	errDomainNameContainsSymbols = errors.New("Team name should be alpha-numeric, no special chars or whitespace is allowed")
	errTeamAlreadyExists         = errors.New("Team name already exists")

	// User
	errUserAlreadyExists         = errors.New("User already exists")
	errUserCreationFailed        = errors.New("User creation failed ")
	errNoTeamIDNotExist          = errors.New("Team ID does not exist")
	errNOTeamIDProvided          = errors.New("No team id provided")
	errVerificationCodeIsEmpty   = errors.New("Verification code is empty")
	errVerificationCodeExpired   = errors.New("Verification code expired")
	errVerificationCodeFailed    = errors.New("Verification code seems to be failed")
	errUsernameOrPasswordIsEmpty = errors.New("User or password is empty ")
	errInvalidEmailOrPassword    = errors.New("Invalid email or password")

	// VCS
	errNoProviderFound         = "No provider found for %s"
	errGetUpdatedFokenFailed   = "Failed to get updated token %s"
	errGettingRepositories     = "Failed to get repositories for %s"
	errVCSAccountAlreadyLinked = errors.New("VCS account already linked")
)

// Common constants
const (
	SLASH     = "/"
	SEMICOLON = ";"
)

// Config ..
type Config struct {
	Github struct {
		Key      string
		Secret   string
		Callback string
	}

	Gitlab struct {
		Key      string
		Secret   string
		Callback string
	}

	Bitbucket struct {
		Key      string
		Secret   string
		Callback string
	}

	DB struct {
		Server    string
		Name      string
		Username  string
		Password  string
		Timeout   int
		Monotonic bool

		// old info
		IdleConnection int
		MaxConnection  int
		Log            bool
		Retry          int
		Reconnect      int
	}
	Key struct {
		VerifyCode string
		Signer     string
		Verifier   string
		Certfile   string
		Keyfile    string
	}
	CSRF struct {
		Key    string
		Secure bool
	}
}

// AppContext ..
type AppContext struct {
	Context  context.Context
	Router   *mux.Router
	Logger   *logrus.Logger
	Signer   interface{}
	Verifier interface{}
	Config   Config

	SecureChain alice.Chain
	PublicChain alice.Chain

	TeamService TeamService
	UserService UserService
	VCSService  VCSService
	RepoService RepoService

	Datasource       Datastore
	TeamDatastore    TeamDatastore
	UserDatastore    UserDatastore
	RepoDatastore    RepoDatastore
	SysconfDatastore SysconfDatastore
}

// VCSSysConf ..(sysconf)
type VCSSysConf struct {
	// common fields for any sys config
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name"`
	Type string        `bson:"type"`

	Key         string `bson:"key"`
	Secret      string `bson:"secret"`
	CallbackURL string `bson:"callback_url"`
	HookURL     string `bson:"hook_url"`
}

// Team ..
type Team struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string        `bson:"name"`
	Display  string        `bson:"display,omitempty"`
	Accounts []VCS         `bson:"accounts,omitempty"`
}

// User ..
type User struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Fullname      string        `bson:"fullname"`
	Username      string        `bson:"username"`
	Email         string        `bson:"email"`
	Password      string        `bson:"password"`
	Locked        int8          `bson:"locked"`
	Active        int8          `bson:"active"`
	BadAttempt    int8          `bson:"bad_attempt"`
	EmailVefified bool          `bson:"email_verified"`
	Scope         []string      `bson:"scope"`
	Team          string        `bson:"team"`
}

// VCS contains the information common amongst most OAuth and OAuth2 providers.
// All of the "raw" datafrom the provider can be found in the `RawData` field.
type VCS struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	Name         string    `json:"name" bson:"name,omitempty"`
	Type         string    `json:"type" bson:"type,omitempty"`
	OwnerType    string    `json:"owner_type" bson:"owner_type,omitempty"`
	AvatarURL    string    `json:"avatar" bson:"avatar,omitempty"`
	AccessCode   string    `json:"-" bson:"access_code,omitempty"`
	AccessToken  string    `json:"-" bson:"access_token,omitempty"`
	RefreshToken string    `json:"-" bson:"refresh_token,omitempty"`
	TokenType    string    `json:"-" bson:"token_type,omitempty"`
	TokenExpiry  time.Time `json:"-" bson:"token_expiry,omitempty"`
}

// Repo ..
// Represents as vcs repositories or projects
type Repo struct {
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Team          string        `json:"-" bson:"team"`
	RepoID        string        `json:"-" bson:"repo_id"`
	VcsID         string        `json:"-" bson:"vcs_id"`
	Name          string        `json:"name" bson:"name,omitempty"`
	Private       bool          `json:"private" bson:"private,omitempty"`
	Link          string        `json:"link" bson:"link,omitempty"`
	Description   string        `json:"description" bson:"description,omitempty"`
	Fork          bool          `json:"fork" bson:"fork,omitempty"`
	DefaultBranch string        `json:"default_branch" bson:"default_branch,omitempty"`
	Language      string        `json:"language" bson:"language,omitempty"`
}
