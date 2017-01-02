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

/*
func (r *sysconfDatastore) CheckExists(name string) (bool, error) {

	var count int
	var err error
	r.ds.Execute(r.cname, func(c *mgo.Collection) {
		count, err = c.Find(bson.M{"name": name}).Count()
	})

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *sysconfDatastore) GetTeamID(name string) (string, error) {

	var err error
	var result Team
	r.ds.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Find(bson.M{"name": name}).One(&result)
	})

	if err != nil {
		return "", err
	}
	return result.ID.String(), nil
}*/

// NewSysconfDatastore ..
func NewSysconfDatastore(ds core.Datastore) SysconfDatastore {
	return &sysconfDatastore{ds: ds, cname: "sysconf"}
}
