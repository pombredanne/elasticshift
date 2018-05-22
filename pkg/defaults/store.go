/*
Copyright 2018 The Elasticshift Authors.
*/
package defaults

import (
	"strings"

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

	FindByReferenceId(id string) (types.Default, error)
	GetDefaultContainerEngine(team string) (string, error)
	UpdateDefaults(referenceID string, fields bson.M) error
}

// NewStore ..
func NewStore(d stypes.Database) Store {
	s := &store{}
	s.Database = d
	s.CollectionName = "default"
	return s
}

func (s *store) FindByReferenceId(id string) (types.Default, error) {

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

func (s *store) GetDefaultContainerEngine(team string) (string, error) {

	result, err := s.FindByReferenceId(team)
	if err != nil {
		return "", err
	}
	return result.ContainerEngineID, nil
}

func (s *store) UpdateDefaults(referenceID string, fields bson.M) error {

	var err error
	s.Execute(func(c *mgo.Collection) {
		err = c.Update(bson.M{"reference_id": referenceID}, bson.M{"$set": fields})
	})
	return err
}
