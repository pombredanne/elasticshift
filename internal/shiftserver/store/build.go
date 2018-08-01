/*
Copyright 2017 The Elasticshift Authors.
*/
package store

import (
	"gitlab.com/conspico/elasticshift/api/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type build struct {
	Store
}

// Build ...
// Store provides build related operation
type Build interface {
	Interface

	FetchBuild(team, repositoryID, branch, id string, status types.BuildStatus) ([]types.Build, error)
	FetchBuildByID(id string) (types.Build, error)
	FetchBuildByRepositoryID(id string) ([]types.Build, error)
	UpdateBuildLog(id bson.ObjectId, log string) error
	UpdateBuildStatus(id bson.ObjectId, s types.BuildStatus) error
	UpdateContainerID(id bson.ObjectId, containerID string) error
}

// NewStore ..
func newBuildStore(d Database) Build {
	s := &build{}
	s.Database = d
	s.CollectionName = "build"
	return s
}

func (s *build) FetchBuild(team, repositoryID, branch, id string, status types.BuildStatus) ([]types.Build, error) {

	q := bson.M{"team": team}
	if repositoryID != "" {
		q["repository_id"] = repositoryID
	}

	if branch != "" {
		q["branch"] = branch
	}

	if status > 0 {
		q["status"] = status
	}

	if id != "" {
		q["_id"] = bson.ObjectIdHex(id)
	}

	var err error
	var result []types.Build
	s.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	return result, err
}

func (s *build) FetchBuildByRepositoryID(id string) ([]types.Build, error) {

	q := bson.M{"repository_id": id}

	var err error
	var result []types.Build
	s.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	return result, err
}

func (s *build) FetchBuildByID(id string) (types.Build, error) {
	var b types.Build
	err := s.FindOne(bson.M{"_id": bson.ObjectIdHex(id)}, &b)
	return b, err
}

func (s *build) UpdateBuildLog(id bson.ObjectId, log string) error {

	var err error
	s.Execute(func(c *mgo.Collection) {
		err = c.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"log": log}})
	})
	return err
}

func (s *build) UpdateBuildStatus(id bson.ObjectId, status types.BuildStatus) error {

	var err error
	s.Execute(func(c *mgo.Collection) {
		err = c.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	})
	return err
}

func (s *build) UpdateContainerID(id bson.ObjectId, containerID string) error {

	var err error
	s.Execute(func(c *mgo.Collection) {
		err = c.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"container_id": containerID}})
	})
	return err
}
