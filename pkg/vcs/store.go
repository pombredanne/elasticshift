/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
	"gitlab.com/conspico/elasticshift/api/types"
	base "gitlab.com/conspico/elasticshift/pkg/store"
	stypes "gitlab.com/conspico/elasticshift/pkg/store/types"
	"gopkg.in/mgo.v2/bson"
)

type store struct {
	base.Store // store
}

//Store related database operations
type Store interface {
	base.Interface

	// VCS Settings
	UpdateVCS(vcs types.VCS) error
}

func NewStore(d stypes.Database) Store {
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
