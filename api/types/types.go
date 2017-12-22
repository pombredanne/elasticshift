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

// VCSSysConf ..(sysconf)
type VCSSysConf struct {
	// common fields for any sys config
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Name string        `bson:"name,omitempty"`
	Kind string        `bson:"kind,omitempty"`

	Key         string `bson:"key,omitempty"`
	Secret      string `bson:"secret,omitempty"`
	CallbackURL string `bson:"callback_url,omitempty"`
	HookURL     string `bson:"hook_url,omitempty"`
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
	// ID           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	ID           string    `json:"id" bson:"id,omitempty"`
	Name         string    `json:"name" bson:"name,omitempty"`
	Kind         int       `json:"kind" bson:"kind,omitempty"`
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
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	RepoID        string        `json:"-" bson:"repo_id"`
	VcsID         string        `json:"-" bson:"vcs_id"`
	Name          string        `json:"name" bson:"name,omitempty"`
	Private       bool          `json:"private" bson:"private,omitempty"`
	Link          string        `json:"link" bson:"link,omitempty"`
	Description   string        `json:"description" bson:"description,omitempty"`
	Fork          bool          `json:"fork" bson:"fork,omitempty"`
	DefaultBranch string        `json:"default_branch" bson:"default_branch,omitempty"`
	Language      string        `json:"language" bson:"language,omitempty"`
	Team          string        `json:"-" bson:"team"`
}

type VCSList struct {
	Nodes []VCS `json:"nodes"`
	Count int   `json:"count"`
}
