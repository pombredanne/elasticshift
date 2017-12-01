/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
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
	GetVCSTypes() ([]VCSSysConf, error)
	SaveVCS(scf *VCSSysConf) error
	Delete(id bson.ObjectId) error
}

// NewStore ..
func NewStore(s core.Store) Store {
	return &store{s, "sysconf"}
}

func (r *store) GetVCSTypes() ([]VCSSysConf, error) {

	result := make([]VCSSysConf, 0)
	err := r.store.FindAll(r.cname, bson.M{"type": vcsType}, &result)
	return result, err
}

func (r *store) SaveVCS(v *VCSSysConf) error {
	return r.store.Insert(r.cname, v)
}

func (r *store) Delete(id bson.ObjectId) error {
	return r.store.Remove(r.cname, id)
}
