/*
Copyright 2017 The Elasticshift Authors.
*/
package sysconf

import (
	"fmt"

	"gitlab.com/conspico/elasticshift/api/types"
	core "gitlab.com/conspico/elasticshift/core/store"
	"gopkg.in/mgo.v2/bson"
)

const (
	vcsKind = "vcs"
)

type store struct {
	store core.Store
	cname string
}

// Store provides system level config
type Store interface {
	GetVCSSysConf() ([]types.VCSSysConf, error)
	GetVCSSysConfByName(name string) (types.VCSSysConf, error)
	SaveVCSSysConf(scf *types.VCSSysConf) error
	Delete(id bson.ObjectId) error
}

// NewStore ..
func NewStore(s core.Store) Store {
	return &store{s, "sysconf"}
}

func (r *store) GetVCSSysConf() ([]types.VCSSysConf, error) {

	result := make([]types.VCSSysConf, 0)
	err := r.store.FindAll(r.cname, bson.M{"kind": vcsKind}, &result)
	return result, err
}

func (r *store) GetVCSSysConfByName(name string) (types.VCSSysConf, error) {

	q := bson.M{}
	q["kind"] = vcsKind
	q["name"] = name

	var result types.VCSSysConf
	fmt.Println("Getting VCS sysconf for :", name)
	err := r.store.FindOne(r.cname, q, &result)
	fmt.Println("Get vcs config:", err)
	return result, err
}

func (r *store) SaveVCSSysConf(v *types.VCSSysConf) error {
	v.Kind = vcsKind
	return r.store.Insert(r.cname, v)
}

func (r *store) Delete(id bson.ObjectId) error {
	return r.store.Remove(r.cname, id)
}
