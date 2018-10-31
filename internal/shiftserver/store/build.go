/*
Copyright 2017 The Elasticshift Authors.
*/
package store

import (
	"fmt"

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

	FetchBuild(team, repositoryID, branch, id string, status []string) ([]types.Build, error)
	FetchBuildByID(id string) (types.Build, error)
	FetchBuildByRepositoryID(id string) ([]types.Build, error)

	UpdateBuildLog(id bson.ObjectId, log string) error
	UpdateBuildStatus(id bson.ObjectId, s string) error
	UpdateContainerID(id bson.ObjectId, containerID string) error

	SaveSubBuild(buildID string, sb *types.SubBuild) error
	UpdateSubBuild(buildID string, sb types.SubBuild) error
	FetchSubBuild(buildID, subBuildID string) (types.SubBuild, error)
}

// NewStore ..
func newBuildStore(d Database) Build {
	s := &build{}
	s.Database = d
	s.CollectionName = "build"
	return s
}

func (s *build) FetchBuild(team, repositoryID, branch, id string, status []string) ([]types.Build, error) {

	q := bson.M{"team": team}
	if repositoryID != "" {
		q["repository_id"] = repositoryID
	}

	if branch != "" {
		q["branch"] = branch
	}

	if statusLen := len(status); statusLen > 0 {
		if statusLen == 1 {
			q["sub_builds.status"] = status[0]
		} else {
			q["sub_builds.status"] = bson.M{"$in": status}
		}
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

func (s *build) UpdateBuildStatus(id bson.ObjectId, status string) error {

	var err error
	s.Execute(func(c *mgo.Collection) {
		err = c.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"sub_builds.$.status": status}})
	})
	return err
}

func (s *build) UpdateContainerID(id bson.ObjectId, containerID string) error {
	return s.UpdateId(id, bson.M{"$set": bson.M{"container_id": containerID}})
}

func (s *build) SaveSubBuild(buildID string, sb *types.SubBuild) error {

	var err error
	s.Execute(func(c *mgo.Collection) {
		err = c.Update(
			bson.M{"_id": bson.ObjectIdHex(buildID)},
			bson.M{"$push": bson.M{"sub_builds": sb}},
		)
	})
	return err
}

func (s *build) UpdateSubBuild(buildID string, sb types.SubBuild) error {

	u := bson.M{}

	if sb.Image != "" {
		u["sub_builds.$.image"] = sb.Image
	}

	if sb.Graph != "" {
		u["sub_builds.$.graph"] = sb.Graph
	}

	if sb.Status != "" {
		u["sub_builds.$.status"] = sb.Status

	}

	if sb.Reason != "" {
		u["sub_builds.$.reason"] = sb.Reason
	}

	if sb.Duration != "" {
		u["sub_builds.$.duration"] = sb.Duration
	}

	var err error
	s.Execute(func(c *mgo.Collection) {
		err = c.Update(
			bson.M{"_id": bson.ObjectIdHex(buildID), "sub_builds.id": sb.ID},
			bson.M{"$set": u})
	})
	return err
}

func (s *build) FetchSubBuild(buildID, subBuildID string) (types.SubBuild, error) {

	var sb types.SubBuild
	var b types.Build
	var err error

	s.Execute(func(c *mgo.Collection) {
		err = c.Find(bson.M{"_id": bson.ObjectIdHex(buildID), "sub_builds.id": subBuildID}).Select(bson.M{"sub_builds.$": 1}).One(&b)
	})

	fmt.Printf("Len of sub_builds : %d", len(b.SubBuilds))

	if len(b.SubBuilds) > 0 {
		sb = b.SubBuilds[0]
	}

	fmt.Printf("Inside FetchSubBuild: ID = %s, Image = %s \n", sb.ID, sb.Image)
	return sb, err
}
