package esh

import (
	"sync"

	"github.com/jinzhu/gorm"
)

type userDatastore struct {
	mtx sync.RWMutex
	db  *gorm.DB
}

func (r *userDatastore) Save(user *User) error {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.db.NewRecord(user)
	r.db.Create(&user)

	return nil
}

func (r *userDatastore) CheckExists(email, teamID string) (bool, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result struct {
		Exist int
	}
	err := r.db.Raw("SELECT 1 as 'exist' FROM USER WHERE email = ? AND TEAM_ID = ? LIMIT 1", email, teamID).Scan(&result).Error

	return result.Exist == 1, err
}

func (r *userDatastore) GetUser(email, teamID string) (User, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result User
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

func (r *userDatastore) FindByName(name string) (User, error) {

	r.mtx.Lock()
	defer r.mtx.Unlock()

	var result User
	err := r.db.Where("username = ?", name).First(&result).Error
	return result, err
}

// NewUserDatastore ..
func NewUserDatastore(db *gorm.DB) UserDatastore {
	return &userDatastore{db: db}
}
