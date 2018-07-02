/*
Copyright 2017 The Elasticshift Authors.
*/
package client

import (
	base "gitlab.com/conspico/elasticshift/internal/shiftserver/store"
	stypes "gitlab.com/conspico/elasticshift/internal/types"
	"gopkg.in/mgo.v2/bson"
)

type store struct {
	base.Store // store
}

// NewStore related database operations
func NewStore(d stypes.Database) Store {
	s := &store{}
	s.Database = d
	s.CollectionName = "client"
	return s
}

// Store related database operations
type Store interface {
	base.Interface

	ExistByName(name string) (bool, error)
}

func (s *store) ExistByName(name string) (bool, error) {
	return s.Exist(bson.M{"name": name})
}
