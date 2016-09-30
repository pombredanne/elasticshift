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

func (r *teamRepository) Save(team *team.Team) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.db.NewRecord(team)
	err := r.db.Create(&team).Error

	return err
}

func (r *teamRepository) CheckExists(name string) (bool, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result struct {
		Exist int
	}
	err := r.db.Raw("SELECT 1 as 'exist' FROM TEAM WHERE name = ? LIMIT 1", name).Scan(&result).Error

	return result.Exist == 1, err
}

func (r *teamRepository) GetTeamID(name string) (string, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result struct {
		ID string
	}
	err := r.db.Raw("SELECT id FROM TEAM WHERE name = ? LIMIT 1", name).Scan(&result).Error
	return result.ID, err
}

func (r *teamRepository) FindByName(name string) (team.Team, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result team.Team
	err := r.db.Where("name = ?", name).First(&result).Error
	return result, err
}

// NewTeam ..
func NewTeam(db *gorm.DB) team.Repository {
	return &teamRepository{db: db}
}
