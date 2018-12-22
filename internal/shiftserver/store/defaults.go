/*
Copyright 2018 The Elasticshift Authors.
*/
package store

import (
	"strings"

	"github.com/elasticshift/elasticshift/api/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type defaults struct {
	Store
}

// Defaults..
// Store provides system level defaults
type Defaults interface {
	Interface

	FindByReferenceId(id string) (types.Default, error)
	GetDefaultContainerEngine(team string) (string, error)
	UpdateDefaults(referenceID string, fields bson.M) error
}

// NewStore ..
func newDefaultsStore(d Database) Defaults {
	s := &defaults{}
	s.Database = d
	s.CollectionName = "default"
	return s
}

func (s *defaults) FindByReferenceId(id string) (types.Default, error) {

	var result types.Default
	err := s.FindOne(bson.M{"reference_id": id}, &result)

	var notfound bool
	if err != nil {
		notfound = strings.EqualFold("not found", err.Error())
	}

	if err != nil && !notfound {
		return types.Default{}, err
	}

	if notfound {
		err = nil
	}
	return result, err
}

func (s *defaults) GetDefaultContainerEngine(team string) (string, error) {

	result, err := s.FindByReferenceId(team)
	if err != nil {
		return "", err
	}
	return result.ContainerEngineID, nil
}

func (s *defaults) UpdateDefaults(referenceID string, fields bson.M) error {

	var err error
	s.Execute(func(c *mgo.Collection) {
		err = c.Update(bson.M{"reference_id": referenceID}, bson.M{"$set": fields})
	})
	return err
}
