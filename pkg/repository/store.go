/*
Copyright 2018 The Elasticshift Authors.
*/
package repository

import (
	"gitlab.com/conspico/elasticshift/api/types"
	base "gitlab.com/conspico/elasticshift/pkg/store"
	stypes "gitlab.com/conspico/elasticshift/pkg/store/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type store struct {
	base.Store // store
}

//Store related database operations
type Store interface {
	base.Interface

	// VCS Settings
	UpdateRepository(repo types.Repository) error
	GetRepositoryByID(id string) (types.Repository, error)
	GetRepository(team, vcsID string) ([]types.Repository, error)
}

// NewStore related database operations
func NewStore(d stypes.Database) Store {
	s := &store{}
	s.Database = d
	s.CollectionName = "repository"
	return s
}

func (s *store) GetRepositoryByID(id string) (types.Repository, error) {

	var result types.Repository
	err := s.FindByID(id, &result)
	return result, err
}

func (s *store) UpdateRepository(vcs types.Repository) error {

	_, err := s.Upsert(bson.M{"_id": vcs.ID}, vcs)
	return err
}

func (s *store) GetRepository(team, vcsID string) ([]types.Repository, error) {

	q := bson.M{"team": team}
	if vcsID != "" {
		q["vcs_id"] = vcsID
	}

	var err error
	var result []types.Repository
	s.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})
	return result, err
}
