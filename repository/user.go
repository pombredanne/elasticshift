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

func (t *userRepository) Save(user *user.User) error {

	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.db.NewRecord(user)
	t.db.Create(&user)

	return nil
}

func (t *userRepository) CheckExists(email string) (bool, error) {

	t.mtx.Lock()
	defer t.mtx.Unlock()

	var result struct {
		Exist int
	}
	err := t.db.Raw("SELECT 1 as 'exist' FROM USERS WHERE email = ? LIMIT 1", email).Scan(&result).Error

	return result.Exist == 1, err
}

func (t *userRepository) FindByName(name string) (user.User, error) {

	t.mtx.Lock()
	defer t.mtx.Unlock()

	var result user.User
	err := t.db.Where("username = ?", name).First(&result).Error
	return result, err
}

// NewUser ..
func NewUser(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}
