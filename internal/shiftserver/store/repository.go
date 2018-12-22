/*
Copyright 2018 The Elasticshift Authors.
*/
package store

import (
	"github.com/elasticshift/elasticshift/api/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type repository struct {
	Store // store
}

//Store related database operations
type Repository interface {
	Interface

	// VCS Settings
	UpdateRepository(repo types.Repository) error
	GetRepositoryByID(id string) (types.Repository, error)
	GetRepository(team, vcsID string) ([]types.Repository, error)
}

// NewStore related database operations
func newRepositoryStore(d Database) Repository {
	s := &repository{}
	s.Database = d
	s.CollectionName = "repository"
	return s
}

func (s *repository) GetRepositoryByID(id string) (types.Repository, error) {

	var result types.Repository
	err := s.FindByID(id, &result)
	return result, err
}

func (s *repository) UpdateRepository(vcs types.Repository) error {

	_, err := s.Upsert(bson.M{"_id": vcs.ID}, vcs)
	return err
}

func (s *repository) GetRepository(team, vcsID string) ([]types.Repository, error) {

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
