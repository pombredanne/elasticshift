package esh

// TeamDatastore provides access a team.
type TeamDatastore interface {
	Save(team *Team) error
	CheckExists(name string) (bool, error)
	GetTeamID(name string) (string, error)
	FindByName(name string) (Team, error)
}

// UserDatastore provides access a user.
type UserDatastore interface {
	Save(user *User) error
	CheckExists(email, teamID string) (bool, error)
	GetUser(email, teamID string) (User, error)
	FindByName(username string) (User, error)
}

// VCSDatastore provides access a user.
type VCSDatastore interface {
	Save(user *VCS) error
	GetVCS(teamID string) ([]VCS, error)
	GetByID(id string) (VCS, error)
	UpdateVCS(old *VCS, updated VCS) error
}

// RepoDatastore provides the repository related datastore func.
type RepoDatastore interface {
	Save(repo *Repo) error
	Update(old Repo, repo Repo) error
	Delete(repo Repo) error
	DeleteIds(ids []string) error
	GetRepos(teamID string) ([]Repo, error)
	GetReposByVCSID(teamID, vcsID string) ([]Repo, error)
}
