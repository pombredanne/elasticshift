/*
Copyright 2017 The Elasticshift Authors.
*/
package store

import (
	"gitlab.com/conspico/elasticshift/api/types"
	"gopkg.in/mgo.v2/bson"
)

const (
	VcsKind       = "vcs"
	GenericKind   = "generic"
	VolumeNfsKind = "nfs"

	DEFAULT_STORAGE = "default-storage"
)

type sysconf struct {
	Store
}

// Store provides system level config
type Sysconf interface {
	Interface

	GetVCSSysConf() ([]types.VCSSysConf, error)
	GetSysConf(kind, name string, result interface{}) error
	GetDefaultStorage() (types.NFSVolumeSysConf, error)
}

// NewStore ..
func newSysconfStore(d Database) Sysconf {
	s := &sysconf{}
	s.Database = d
	s.CollectionName = "sysconf"
	return s
}

func (r *sysconf) GetVCSSysConf() ([]types.VCSSysConf, error) {

	result := make([]types.VCSSysConf, 0)
	err := r.FindAll(bson.M{"kind": VcsKind}, &result)
	return result, err
}

func (r *sysconf) GetSysConf(kind, name string, result interface{}) error {

	q := bson.M{}
	q["kind"] = kind
	q["name"] = name

	return r.FindOne(q, result)
}

func (r *sysconf) GetDefaultStorage() (types.NFSVolumeSysConf, error) {

	// find the system storage
	var result types.NFSVolumeSysConf
	err := r.GetSysConf(VolumeNfsKind, DEFAULT_STORAGE, &result)
	return result, err
}
