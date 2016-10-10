package repository

import (
	"sync"

	"gitlab.com/conspico/esh/vcs"

	"github.com/jinzhu/gorm"
)

type vcsRepository struct {
	mtx sync.RWMutex
	db  *gorm.DB
}

func (r *vcsRepository) Save(v *vcs.VCS) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.db.NewRecord(v)
	err := r.db.Create(&v).Error

	return err
}

func (r *vcsRepository) GetVCS(teamID string) ([]vcs.VCS, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result []vcs.VCS
	err := r.db.Raw(`SELECT id, 
							name, 
							type, 
							avatar_url,
							updated_dt
				     FROM VCS WHERE TEAM_ID = ? LIMIT 1`, teamID).Scan(&result).Error
	return result, err
}

// NewVCS ..
func NewVCS(db *gorm.DB) vcs.Repository {
	return &vcsRepository{db: db}
}
