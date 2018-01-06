/*
Copyright 2017 The Elasticshift Authors.
*/
package build

import (
	"gitlab.com/conspico/elasticshift/api/types"
	core "gitlab.com/conspico/elasticshift/core/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type store struct {
	store core.Store
	cname string
}

// Store provides system level config
type Store interface {
	SaveBuild(b *types.Build) error
	FetchBuild(team, repository_id, branch string, status types.BuildStatus) ([]types.Build, error)
	FetchBuildByID(id string) (types.Build, error)
	UpdateBuildLog(id bson.ObjectId, log string) error
}

// NewStore ..
func NewStore(s core.Store) Store {
	return &store{s, "build"}
}

func (s *store) SaveBuild(b *types.Build) error {
	return s.store.Insert(s.cname, b)
}

func (s *store) FetchBuild(team, repository_id, branch string, status types.BuildStatus) ([]types.Build, error) {

	q := bson.M{"team": team}
	if repository_id != "" {
		q["repository_id"] = repository_id
	}

	if branch != "" {
		q["branch"] = branch
	}

	if status > 0 {
		q["status"] = status
	}

	var err error
	var result []types.Build
	s.store.Execute(s.cname, func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	return result, err
}

func (s *store) FetchBuildByID(id string) (types.Build, error) {
	var b types.Build
	err := s.store.FindOne(s.cname, bson.M{"_id": id}, &b)
	return b, err
}

func (s *store) UpdateBuildLog(id bson.ObjectId, log string) error {

	var err error
	s.store.Execute(s.cname, func(c *mgo.Collection) {
		err = c.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"log": log}})
	})
	return err
}
