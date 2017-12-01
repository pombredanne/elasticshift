/*
Copyright 2017 The Elasticshift Authors.
*/
package team

import (
	"gitlab.com/conspico/elasticshift/api/types"
	core "gitlab.com/conspico/elasticshift/core/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type store struct {
	store core.Store // store
	cname string     // collection name
}

//Store related database operations
type Store interface {
	//Team
	Save(team *types.Team) error
	CheckExists(name string) (bool, error)
	GetTeam(name string) (types.Team, error)

	// VCS Settings
	SaveVCS(team string, vcs *types.VCS) error
	UpdateVCS(team string, vcs types.VCS) error
	GetVCSByID(team, id string) (types.VCS, error)
}

// NewStore related database operations
func NewStore(s core.Store) Store {
	return &store{s, "team"}
}

func (r *store) Save(team *types.Team) error {
	return r.store.Insert(r.cname, team)
}

func (r *store) CheckExists(name string) (bool, error) {

	var count int
	var err error
	r.store.Execute(r.cname, func(c *mgo.Collection) {
		count, err = c.Find(bson.M{"name": name}).Count()
	})

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *store) GetTeam(name string) (types.Team, error) {

	var err error
	var result types.Team
	r.store.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Find(bson.M{"name": name}).One(&result)
	})
	return result, err
}

func (r *store) SaveVCS(team string, vcs *types.VCS) error {

	var err error
	r.store.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Update(
			bson.M{"name": team},
			bson.M{"$push": bson.M{"accounts": vcs}},
		)
	})
	return err
}

func (r *store) GetVCSByID(team, id string) (types.VCS, error) {

	var t types.Team
	var err error
	r.store.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Find(bson.M{"name": team, "accounts._id": id}).Select(bson.M{"accounts.$": 1}).One(&t)
	})

	if len(t.Accounts) == 0 {
		return types.VCS{}, err
	}

	return t.Accounts[0], err
}

func (r *store) UpdateVCS(team string, vcs types.VCS) error {

	var err error
	r.store.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Update(bson.M{"name": team, "accounts._id": vcs.ID},
			bson.M{"$set": bson.M{"accounts.$.access_token": vcs.AccessToken,
				"accounts.$.access_code":   vcs.AccessCode,
				"accounts.$.refresh_token": vcs.RefreshToken,
				"accounts.$.owner_type":    vcs.OwnerType,
				"accounts.$.token_expiry":  vcs.TokenExpiry}})
	})
	return err
}
