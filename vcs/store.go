/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
	"gitlab.com/conspico/elasticshift/api/types"
	core "gitlab.com/conspico/elasticshift/core/store"
	"gopkg.in/mgo.v2/bson"
)

type store struct {
	core.Store // store
}

//Store related database operations
type Store interface {
	core.Core

	// VCS Settings
	UpdateVCS(vcs types.VCS) error
}

func NewStore(d core.Database) Store {
	s := &store{}
	s.Database = d
	s.CollectionName = "vcs"
	return s
}

func (s *store) GetVCSByID(id string) (types.VCS, error) {

	var result types.VCS
	err := s.FindByID(id, &result)
	return result, err
}

func (s *store) UpdateVCS(vcs types.VCS) error {
	_, err := s.Upsert(bson.M{"_id": vcs.ID}, vcs)
	return err
}
