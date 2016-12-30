// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type teamDatastore struct {
	ds    Datastore
	cname string
}

func (r *teamDatastore) Save(team *Team) error {
	return r.ds.Insert(r.cname, team)
}

func (r *teamDatastore) CheckExists(name string) (bool, error) {

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

func (r *teamDatastore) GetTeam(name string) (Team, error) {

	var err error
	var result Team
	r.ds.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Find(bson.M{"name": name}).One(&result)
	})
	return result, err
}

func (r *teamDatastore) SaveVCS(team string, vcs *VCS) error {

	var err error
	r.ds.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Update(
			bson.M{"name": team},
			bson.M{"$push": bson.M{"accounts": vcs}},
		)
	})
	return err
}

func (r *teamDatastore) GetVCSByID(team, id string) (VCS, error) {

	var t Team
	var err error
	r.ds.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Find(bson.M{"name": team, "accounts._id": id}).Select(bson.M{"accounts.$": 1}).One(&t)
	})

	if len(t.Accounts) == 0 {
		return VCS{}, err
	}

	return t.Accounts[0], err
}

func (r *teamDatastore) UpdateVCS(team string, vcs VCS) error {

	var err error
	r.ds.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Update(bson.M{"name": team, "accounts._id": vcs.ID},
			bson.M{"$set": bson.M{"accounts.$.access_token": vcs.AccessToken,
				"accounts.$.access_code":   vcs.AccessCode,
				"accounts.$.refresh_token": vcs.RefreshToken,
				"accounts.$.owner_type":    vcs.OwnerType,
				"accounts.$.token_expiry":  vcs.TokenExpiry}})
	})
	return err
}

// NewTeamDatastore ..
func NewTeamDatastore(ds Datastore) TeamDatastore {
	return &teamDatastore{ds: ds, cname: "teams"}
}
