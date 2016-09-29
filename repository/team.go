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
	err := t.db.Create(&team).Error

	return err
}

func (t *teamRepository) CheckExists(name string) (bool, error) {

	t.mtx.Lock()
	defer t.mtx.Unlock()

	var result struct {
		Exist int
	}
	err := t.db.Raw("SELECT 1 as 'exist' FROM TEAM WHERE name = ? LIMIT 1", name).Scan(&result).Error

	return result.Exist == 1, err
}

func (t *teamRepository) GetTeamID(name string) (string, error) {

	t.mtx.Lock()
	defer t.mtx.Unlock()

	var result struct {
		ID string
	}
	err := t.db.Raw("SELECT id FROM TEAM WHERE name = ? LIMIT 1", name).Scan(&result).Error
	return result.ID, err
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
