// Package store ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package store

import "gopkg.in/mgo.v2/bson"

const (
	vcsType = "vcs"
)

type sysconfStore struct {
	store Store
	cname string
}

// SysconfStore provides system level config
type SysconfStore interface {
	GetVCSTypes() ([]VCSSysConf, error)
	SaveVCS(scf *VCSSysConf) error
	Delete(id bson.ObjectId) error
}

func (r *sysconfStore) GetVCSTypes() ([]VCSSysConf, error) {

	result := make([]VCSSysConf, 0)
	err := r.store.FindAll(r.cname, bson.M{"type": vcsType}, &result)
	return result, err
}

func (r *sysconfStore) SaveVCS(v *VCSSysConf) error {
	return r.store.Insert(r.cname, v)
}

func (r *sysconfStore) Delete(id bson.ObjectId) error {
	return r.store.Remove(r.cname, id)
}

// NewSysconfStore ..
func NewSysconfStore(s Store) SysconfStore {
	return &sysconfStore{s, "sysconf"}
}
