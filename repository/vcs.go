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

func (r *vcsRepository) Save(v *vcs.User) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.db.NewRecord(v)
	r.db.Create(&v)

	return nil
}

// NewVCS ..
func NewVCS(db *gorm.DB) vcs.Repository {
	return &vcsRepository{db: db}
}
