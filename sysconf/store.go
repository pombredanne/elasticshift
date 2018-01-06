/*
Copyright 2017 The Elasticshift Authors.
*/
package sysconf

import (
	"gitlab.com/conspico/elasticshift/api/types"
	core "gitlab.com/conspico/elasticshift/core/store"
	"gopkg.in/mgo.v2/bson"
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

	SaveGenericSysConf(scf *types.GenericSysConf) error
	GetGenericSysConfByName(name string) (types.GenericSysConf, error)

	SaveSysConf(obj interface{}) error
	GetSysConf(kind, name string, result interface{}) error
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
	err := r.store.FindOne(r.cname, q, &result)
	return result, err
}
func (r *store) SaveGenericSysConf(g *types.GenericSysConf) error {
	g.Kind = genericKind
	return r.store.Insert(r.cname, g)
}

func (r *store) SaveVCSSysConf(v *types.VCSSysConf) error {
	v.Kind = vcsKind
	return r.store.Insert(r.cname, v)
}

func (r *store) Delete(id bson.ObjectId) error {
	return r.store.Remove(r.cname, id)
}

func (r *store) GetGenericSysConfByName(name string) (types.GenericSysConf, error) {

	q := bson.M{}
	q["kind"] = genericKind
	q["name"] = name

	var result types.GenericSysConf
	err := r.store.FindOne(r.cname, q, &result)
	return result, err
}

func (r *store) SaveSysConf(obj interface{}) error {
	return r.store.Insert(r.cname, &obj)
}

func (r *store) GetSysConf(kind, name string, result interface{}) error {

	q := bson.M{}
	q["kind"] = kind
	q["name"] = name

	return r.store.FindOne(r.cname, q, result)
}
