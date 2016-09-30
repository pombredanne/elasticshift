package repository

import (
	"sync"

	"gitlab.com/conspico/esh/user"

	"github.com/jinzhu/gorm"
)

type userRepository struct {
	mtx sync.RWMutex
	db  *gorm.DB
}

func (r *userRepository) Save(user *user.User) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.db.NewRecord(user)
	r.db.Create(&user)

	return nil
}

func (r *userRepository) CheckExists(email, teamID string) (bool, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result struct {
		Exist int
	}
	err := r.db.Raw("SELECT 1 as 'exist' FROM USER WHERE email = ? AND TEAM_ID = ? LIMIT 1", email, teamID).Scan(&result).Error

	return result.Exist == 1, err
}

func (r *userRepository) GetUser(email, teamID string) (user.User, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result user.User
	err := r.db.Raw(`SELECT id, 
							team_id, 
							fullname, 
							username,
							email,
							password,
							locked,
							active,
							bad_attempt
				     FROM USER WHERE email = ? AND TEAM_ID = ? LIMIT 1`, email, teamID).Scan(&result).Error
	return result, err
}

func (r *userRepository) FindByName(name string) (user.User, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result user.User
	err := r.db.Where("username = ?", name).First(&result).Error
	return result, err
}

// NewUser ..
func NewUser(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}
