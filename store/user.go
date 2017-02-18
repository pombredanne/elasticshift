package store

import "gopkg.in/mgo.v2/bson"

type userStore struct {
	store Store  // store
	cname string // collection name
}

// NewUserStore related database operations
func NewUserStore(s Store) UserStore {
	return &userStore{s, "users"}
}

// UserStore related database operations
type UserStore interface {
	Insert(u *User) error
	GetUser(email, teamname string) (User, error)
}

func (s *userStore) Insert(u *User) error {
	return s.store.Insert(s.cname, u)
}

func (s *userStore) GetUser(email, teamname string) (User, error) {

	var result User
	err := s.store.FindOne(s.cname, bson.M{"email": email, "team": teamname}, &result)
	return result, err
}
