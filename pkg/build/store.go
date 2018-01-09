/*
Copyright 2017 The Elasticshift Authors.
*/
package build

import (
	"gitlab.com/conspico/elasticshift/api/types"
	base "gitlab.com/conspico/elasticshift/pkg/store"
	stypes "gitlab.com/conspico/elasticshift/pkg/store/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type store struct {
	base.Store
}

// Store provides system level config
type Store interface {
	base.Interface

	FetchBuild(team, repository_id, branch string, status types.BuildStatus) ([]types.Build, error)
	FetchBuildByID(id string) (types.Build, error)
	UpdateBuildLog(id bson.ObjectId, log string) error
	UpdateBuildStatus(id bson.ObjectId, s types.BuildStatus) error
}

// NewStore ..
func NewStore(d stypes.Database) Store {
	s := &store{}
	s.Database = d
	s.CollectionName = "build"
	return s
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
	s.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	return result, err
}

func (s *store) FetchBuildByID(id string) (types.Build, error) {
	var b types.Build
	err := s.FindOne(bson.M{"_id": id}, &b)
	return b, err
}

func (s *store) UpdateBuildLog(id bson.ObjectId, log string) error {

	var err error
	s.Execute(func(c *mgo.Collection) {
		err = c.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"log": log}})
	})
	return err
}

func (s *store) UpdateBuildStatus(id bson.ObjectId, status types.BuildStatus) error {

	var err error
	s.Execute(func(c *mgo.Collection) {
		err = c.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	})
	return err
}
