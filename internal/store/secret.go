/*
Copyright 2018 The Elasticshift Authors.
*/
package store

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type secret struct {
	Store
}

// Store provides system level config
type Secret interface {
	Interface

	UpdateSecret(id bson.ObjectId, fields bson.M) error
}

// NewStore ..
func newSecretStore(d Database) Secret {
	s := &secret{}
	s.Database = d
	s.CollectionName = "secret"
	return s
}

func (s *secret) UpdateSecret(id bson.ObjectId, fields bson.M) error {

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
