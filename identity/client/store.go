/*
Copyright 2017 The Elasticshift Authors.
*/
package client

import (
	core "gitlab.com/conspico/elasticshift/core/store"
	"gopkg.in/mgo.v2/bson"
)

type store struct {
	core.Store // store
}

// NewStore related database operations
func NewStore(d core.Database) Store {
	s := &store{}
	s.Database = d
	s.CollectionName = "client"
	return s
}

// Store related database operations
type Store interface {
	core.Core

	ExistByName(name string) (bool, error)
}

func (s *store) ExistByName(name string) (bool, error) {
	return s.Exist(bson.M{"name": name})
}
