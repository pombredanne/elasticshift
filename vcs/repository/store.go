/*
Copyright 2017 The Elasticshift Authors.
*/
package repository

import (
	"gitlab.com/conspico/elasticshift/api/types"
	core "gitlab.com/conspico/elasticshift/core/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type store struct {
	store core.Store // store
	cname string     // collection name
}

//Store related database operations
type Store interface {

	// VCS Settings
	SaveRepository(repo *types.Repository) error
	UpdateRepository(repo types.Repository) error
	GetRepositoryByID(id string) (types.Repository, error)
	GetRepository(team string) ([]types.Repository, error)
}

// NewStore related database operations
func NewStore(s core.Store) Store {
	return &store{store: s, cname: "repository"}
}

func (s *store) SaveRepository(repo *types.Repository) error {
	return s.store.Insert(s.cname, repo)
}

func (s *store) GetRepositoryByID(id string) (types.Repository, error) {

	var result types.Repository
	err := s.store.FindOne(s.cname, bson.M{"_id": bson.ObjectIdHex(id)}, &result)
	return result, err
}

func (s *store) UpdateRepository(vcs types.Repository) error {

	_, err := s.store.Upsert(s.cname, bson.M{"_id": vcs.ID}, vcs)
	return err
}

func (s *store) GetRepository(team string) ([]types.Repository, error) {

	var err error
	var result []types.Repository
	s.store.Execute(s.cname, func(c *mgo.Collection) {
		err = c.Find(bson.M{"team": team}).All(&result)
	})
	return result, err
}
