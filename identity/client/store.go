/*
Copyright 2017 The Elasticshift Authors.
*/
package client

import (
	"gitlab.com/conspico/elasticshift/api/types"
	core "gitlab.com/conspico/elasticshift/core/store"
	"gopkg.in/mgo.v2/bson"
)

type store struct {
	store core.Store // store
	cname string     // collection name
}

// NewStore related database operations
func NewStore(s core.Store) Store {
	return &store{s, "client"}
}

// Store related database operations
type Store interface {
	Insert(c *types.Client) error
	Exist(name string) (bool, error)
	Delete(id string) error
	FindOne(id string) (types.Client, error)
}

func (s *store) Insert(c *types.Client) error {
	return s.store.Insert(s.cname, c)
}

func (s *store) Delete(id string) error {
	return s.store.Remove(s.cname, id)
}

func (s *store) Exist(name string) (bool, error) {
	return s.store.Exist(s.cname, bson.M{"name": name})
}

func (s *store) FindOne(id string) (types.Client, error) {

	var c types.Client
	err := s.store.FindOne(s.cname, bson.M{"id": id}, &c)
	return c, err
}
