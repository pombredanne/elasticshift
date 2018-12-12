/*
Copyright 2018 The Elasticshift Authors.
*/
package store

import (
	"gitlab.com/conspico/elasticshift/api/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type integration struct {
	Store
}

// Store provides system level config
type Integration interface {
	Interface

	// Kube config
	GetKubeConfig(team string) (types.KubeConfig, error)
	SaveKubeConfig(team string, f types.KubeConfig) error

	UpdateWorkerPath(id bson.ObjectId, path string) error
}

// NewStore ..
func newIntegrationStore(d Database) Integration {
	s := &integration{}
	s.Database = d
	s.CollectionName = "integration"
	return s
}

func (s *integration) SaveKubeConfig(team string, f types.KubeConfig) error {

	var err error
	s.Execute(func(c *mgo.Collection) {
		err = c.Update(
			bson.M{"name": team},
			bson.M{"$push": bson.M{"kube_config": f}},
		)
	})
	return err
}

func (r *integration) GetKubeConfig(team string) (types.KubeConfig, error) {

	var err error
	var t types.Team
	r.Execute(func(c *mgo.Collection) {
		err = c.Find(bson.M{"name": team}).One(&t)
	})
	return t.KubeConfig, err
}

func (r *integration) UpdateWorkerPath(id bson.ObjectId, path string) error {

	var err error
	r.Execute(func(c *mgo.Collection) {
		err = c.Update(bson.M{"_id": id},
			bson.M{"$set": bson.M{
				"worker_path": path,
			}})
	})
	return err
}
