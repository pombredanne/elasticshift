// Package esh ...
// Author: Ghazni Nattarshah
// Date: Nov 22, 2016
package esh

import (
	"gopkg.in/mgo.v2/bson"
)

// TeamDatastore provides access a team.
type TeamDatastore interface {

	//Team
	Save(team *Team) error
	CheckExists(name string) (bool, error)
	GetTeam(name string) (Team, error)

	// VCS Settings
	SaveVCS(team string, vcs *VCS) error
	UpdateVCS(team string, vcs VCS) error
	GetVCSByID(team, id string) (VCS, error)
}

// UserDatastore provides access a user.
type UserDatastore interface {
	Save(user *User) error
	CheckExists(email, teamname string) (bool, error)
	GetUser(email, teamname string) (User, error)
}

// RepoDatastore provides the repository related datastore func.
type RepoDatastore interface {
	Save(repo *Repo) error
	Update(repo Repo) error
	Delete(repo Repo) error
	DeleteIds(ids []bson.ObjectId) error
	GetRepos(teamID string) ([]Repo, error)
	GetReposByVCSID(team, vcsID string) ([]Repo, error)
}

// SysconfDatastore provides system level config
type SysconfDatastore interface {
	GetVCSTypes() ([]VCSSysConf, error)
	SaveVCS(scf *VCSSysConf) error
	Delete(id bson.ObjectId) error
}
