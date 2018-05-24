/*
Copyright 2018 The Elasticshift Authors.
*/
package secret

import (
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

	UpdateSecret(id bson.ObjectId, fields bson.M) error
}

// NewStore ..
func NewStore(d stypes.Database) Store {
	s := &store{}
	s.Database = d
	s.CollectionName = "secret"
	return s
}

func (s *store) UpdateSecret(id bson.ObjectId, fields bson.M) error {

	var err error
	s.Execute(func(c *mgo.Collection) {
		err = c.Update(bson.M{"_id": id}, bson.M{"$set": fields})
	})
	return err
}

// func SaveSecret(sec types.Secret) error {

// }

// func UpdateSecret(sec types.Secret) error

// }
