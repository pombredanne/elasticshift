package esh

import (
	"sync"

	"github.com/jinzhu/gorm"
)

type repoDatastore struct {
	mtx sync.RWMutex
	db  *gorm.DB
}

func (r *repoDatastore) Save(repo *Repo) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.db.NewRecord(repo)
	err := r.db.Create(&repo).Error

	return err
}

func (r *repoDatastore) GetReposByVCSID(teamID, vcsID string) ([]Repo, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	result := []Repo{}
	err := r.db.Raw(`SELECT id,
							team_id,
							vcs_id,
							repo_id,
							name,
							private,
							link,
							description,
							fork,
							default_branch,
							language
						FROM REPO WHERE team_id = ? and vcs_id = ?`, teamID, vcsID).Scan(&result).Error
	return result, err
}

func (r *repoDatastore) Update(old Repo, repo Repo) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.db.Model(&old).Updates(repo).Error
}

func (r *repoDatastore) Delete(repo Repo) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.db.Delete(&repo).Error
}

func (r *repoDatastore) DeleteIds(ids []string) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	return r.db.Delete(Repo{}, "ID IN (?)", ids).Error
}

func (r *repoDatastore) GetRepos(teamID string) ([]Repo, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result []Repo
	err := r.db.Raw(`SELECT id,
							team_id,
							vcs_id,
							repo_id,
							name,
							private,
							link,
							description,
							fork,
							default_branch,
							language
					FROM REPO WHERE team_id = ?`, teamID).Scan(&result).Error
	return result, err
}

// NewRepoDatastore ..
func NewRepoDatastore(db *gorm.DB) RepoDatastore {
	return &repoDatastore{db: db}
}
