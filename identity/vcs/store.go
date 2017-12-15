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
	store core.Store // store
	cname string     // collection name
}

//Store related database operations
type Store interface {

	// VCS Settings
	SaveVCS(vcs *types.VCS) error
	UpdateVCS(vcs types.VCS) error
	GetVCSByID(id string) (types.VCS, error)
}

// NewStore related database operations
func NewStore(s core.Store) Store {
	return &store{store: s, cname: "vcs"}
}

func (s *store) SaveVCS(vcs *types.VCS) error {
	return s.store.Insert(s.cname, vcs)
}

func (s *store) GetVCSByID(id string) (types.VCS, error) {

	var result types.VCS
	err := s.store.FindOne(s.cname, bson.M{"_id": bson.ObjectIdHex(id)}, &result)
	return result, err
}

func (s *store) UpdateVCS(vcs types.VCS) error {

	var err error
	return err
}
