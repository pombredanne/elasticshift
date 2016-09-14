package repository

import (
	"sync"

	"gitlab.com/conspico/esh/team"

	"github.com/jinzhu/gorm"
)

type teamRepository struct {
	mtx sync.RWMutex
	db  *gorm.DB
}

func (t *teamRepository) Save(team *team.Team) error {

	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.db.NewRecord(team)
	t.db.Create(&team)

	return nil
}

// NewTeam ..
func NewTeam(db *gorm.DB) team.Repository {
	return &teamRepository{db: db}
}

// func (t *teamRepository) Find(puuid string) team.Team {

// }
