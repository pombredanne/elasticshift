/*
Copyright 2017 The Elasticshift Authors.
*/
package sysconf

import (
	"gitlab.com/conspico/elasticshift/api/types"
	base "gitlab.com/conspico/elasticshift/pkg/store"
	stypes "gitlab.com/conspico/elasticshift/pkg/store/types"
	"gopkg.in/mgo.v2/bson"
)

type store struct {
	base.Store
}

// Store provides system level config
type Store interface {
	base.Interface

	GetVCSSysConf() ([]types.VCSSysConf, error)

	GetSysConf(kind, name string, result interface{}) error
}

// NewStore ..
func NewStore(d stypes.Database) Store {
	s := &store{}
	s.Database = d
	s.CollectionName = "sysconf"
	return s
}

func (r *store) GetVCSSysConf() ([]types.VCSSysConf, error) {

	result := make([]types.VCSSysConf, 0)
	err := r.FindAll(bson.M{"kind": VcsKind}, &result)
	return result, err
}

func (r *store) GetSysConf(kind, name string, result interface{}) error {

	q := bson.M{}
	q["kind"] = kind
	q["name"] = name

	return r.FindOne(q, result)
}
