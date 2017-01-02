// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

import (
	"gopkg.in/mgo.v2/bson"
	"gitlab.com/conspico/esh/core"
)

const (
	vcsType = "vcs"
)

type sysconfDatastore struct {
	ds    core.Datastore
	cname string
}

func (r *sysconfDatastore) GetVCSTypes() ([]VCSSysConf, error) {

	result := make([]VCSSysConf, 0)
	err := r.ds.FindAll(r.cname, bson.M{"type": vcsType}, &result)
	return result, err
}

func (r *sysconfDatastore) SaveVCS(v *VCSSysConf) error {
	return r.ds.Insert(r.cname, v)
}

func (r *sysconfDatastore) Delete(id bson.ObjectId) error {
	return r.ds.Remove(r.cname, id)
}

// NewSysconfDatastore ..
func NewSysconfDatastore(ds core.Datastore) SysconfDatastore {
	return &sysconfDatastore{ds: ds, cname: "sysconf"}
}
