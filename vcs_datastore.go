package esh

import (
	"sync"

	"github.com/jinzhu/gorm"
)

type vcsDatastore struct {
	mtx sync.RWMutex
	db  *gorm.DB
}

func (r *vcsDatastore) Save(v *VCS) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.db.NewRecord(v)
	err := r.db.Create(&v).Error

	return err
}

func (r *vcsDatastore) GetVCS(teamID string) ([]VCS, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result []VCS
	err := r.db.Raw(`SELECT id, 
							name, 
							type, 
							avatar_url,
							updated_dt
				     FROM VCS WHERE TEAM_ID = ? LIMIT 1`, teamID).Scan(&result).Error
	return result, err
}

func (r *vcsDatastore) GetByID(id string) (VCS, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result VCS
	err := r.db.Raw(`SELECT id, 
							name, 
							type, 
							avatar_url,
							access_token,
							refresh_token,
							token_expiry,
							token_type
				     FROM VCS WHERE ID = ? LIMIT 1`, id).Scan(&result).Error
	return result, err
}

func (r *vcsDatastore) UpdateVCS(old *VCS, updated VCS) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	err := r.db.Model(old).Updates(updated).Error
	return err
}

// NewVCSDatastore ..
func NewVCSDatastore(db *gorm.DB) VCSDatastore {
	return &vcsDatastore{db: db}
}
