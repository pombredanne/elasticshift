// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

import (
	"gitlab.com/conspico/esh/core"
	"gopkg.in/mgo.v2/bson"
)

type userDatastore struct {
	ds    core.Datastore
	cname string
}

func (r *userDatastore) Save(user *User) error {
	return r.ds.Insert(r.cname, user)
}

func (r *userDatastore) CheckExists(email, teamname string) (bool, error) {
	return r.ds.Exist(r.cname, bson.M{"email": email, "team": teamname})
}

func (r *userDatastore) GetUser(email, teamname string) (User, error) {

	var result User
	err := r.ds.FindOne(r.cname, bson.M{"email": email, "team": teamname}, &result)
	return result, err
}

// NewUserDatastore ..
func NewUserDatastore(ds core.Datastore) UserDatastore {
	return &userDatastore{ds: ds, cname: "oauth_users"}
}
