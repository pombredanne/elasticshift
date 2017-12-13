/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
	"gitlab.com/conspico/elasticshift/api/types"
	core "gitlab.com/conspico/elasticshift/core/store"
	"gopkg.in/mgo.v2/bson"
)

const (
	vcsType = "vcs"
)

type store struct {
	store core.Store
	cname string
}

// Store provides system level config
type Store interface {
	GetVCSTypes() ([]types.VCSSysConf, error)
	SaveVCS(scf *types.VCSSysConf) error
	Delete(id bson.ObjectId) error
}

// NewStore ..
func NewStore(s core.Store) Store {
	return &store{s, "sysconf"}
}

func (r *store) GetVCSTypes() ([]types.VCSSysConf, error) {

	result := make([]types.VCSSysConf, 0)
	err := r.store.FindAll(r.cname, bson.M{"type": vcsType}, &result)
	return result, err
}

func (r *store) SaveVCS(v *types.VCSSysConf) error {
	return r.store.Insert(r.cname, v)
}

func (r *store) Delete(id bson.ObjectId) error {
	return r.store.Remove(r.cname, id)
}
