package store

import "gopkg.in/mgo.v2/bson"

type clientStore struct {
	store Store  // store
	cname string // collection name
}

// NewClientStore related database operations
func NewClientStore(s Store) ClientStore {
	return &clientStore{s, "client"}
}

// ClientStore related database operations
type ClientStore interface {
	Insert(c *Client) error
	Exist(name string) (bool, error)
	Delete(id string) error
	FindOne(id string) (Client, error)
}

func (s *clientStore) Insert(c *Client) error {
	return s.store.Insert(s.cname, c)
}

func (s *clientStore) Delete(id string) error {
	return s.store.Remove(s.cname, id)
}

func (s *clientStore) Exist(name string) (bool, error) {
	return s.store.Exist(s.cname, bson.M{"name": name})
}

func (s *clientStore) FindOne(id string) (Client, error) {

	var c Client
	err := s.store.FindOne(s.cname, bson.M{"id": id}, &c)
	return c, err
}
