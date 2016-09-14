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

func (t *teamRepository) CheckExists(name string) (bool, error) {

	t.mtx.Lock()
	defer t.mtx.Unlock()

	var result struct {
		Exist int
	}
	err := t.db.Raw("SELECT 1 as 'exist' FROM TEAMS WHERE name = ? LIMIT 1", name).Scan(&result).Error

	return result.Exist == 1, err
}

func (t *teamRepository) FindByName(name string) (team.Team, error) {

	t.mtx.Lock()
	defer t.mtx.Unlock()

	var result team.Team
	err := t.db.Where("name = ?", name).First(&result).Error
	return result, err
}

// NewTeam ..
func NewTeam(db *gorm.DB) team.Repository {
	return &teamRepository{db: db}
}
