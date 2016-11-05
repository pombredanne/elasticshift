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

func (r *repoDatastore) GetReposByVCSID(id string) ([]Repo, error) {

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
					     FROM REPO WHERE vcs_id = ?`, id).Scan(&result).Error
	//err := r.db.Find(&result, "vcs_id = ?", id).Error
	return result, err
}

func (r *repoDatastore) Update(old Repo, repo Repo) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	err := r.db.Model(&old).Updates(repo).Error
	return err
}

func (r *repoDatastore) Delete(repo Repo) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	err := r.db.Delete(&repo).Error
	return err
}

// NewRepoDatastore ..
func NewRepoDatastore(db *gorm.DB) RepoDatastore {
	return &repoDatastore{db: db}
}
