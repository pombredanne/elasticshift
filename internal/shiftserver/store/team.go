/*
Copyright 2017 The Elasticshift Authors.
*/
package store

import (
	"fmt"

	"gitlab.com/conspico/elasticshift/api/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type team struct {
	Store // store
}

//Store related database operations
type Team interface {
	Interface

	//Team
	CheckExists(name string) (bool, error)
	GetTeam(id, name string) (types.Team, error)
	GetTeams(limit, offset int) ([]types.Team, error)

	// VCS Settings
	GetVCS(team string) ([]types.VCS, error)
	SaveVCS(team string, vcs *types.VCS) error
	UpdateVCS(team string, vcs types.VCS) error
	GetVCSByID(team, id string) (types.VCS, error)
	GetVCSByName(team, name, source string) (*types.VCS, error)
}

// NewStore related database operations
func newTeamStore(d Database) Team {
	s := &team{}
	s.Database = d
	s.CollectionName = "team"
	return s
}

func (r *team) GetTeams(limit, offset int) ([]types.Team, error) {

	list := make([]types.Team, 1)

	// TODO fetch list of teams

	return list, nil
}

func (r *team) CheckExists(name string) (bool, error) {

	var count int
	var err error
	r.Execute(func(c *mgo.Collection) {
		count, err = c.Find(bson.M{"name": name}).Count()
	})

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *team) GetTeam(id, name string) (types.Team, error) {

	q := bson.M{}
	if id != "" {
		q["_id"] = bson.ObjectIdHex(id)
	}

	if name != "" {
		q["name"] = name
	}

	var err error
	var result types.Team
	r.Execute(func(c *mgo.Collection) {
		err = c.Find(q).One(&result)
	})
	return result, err
}

func (r *team) SaveVCS(team string, vcs *types.VCS) error {

	var err error
	r.Execute(func(c *mgo.Collection) {
		err = c.Update(
			bson.M{"name": team},
			bson.M{"$push": bson.M{"accounts": vcs}},
		)
	})
	return err
}

func (r *team) GetVCS(team string) ([]types.VCS, error) {

	var err error
	var t types.Team
	r.Execute(func(c *mgo.Collection) {
		err = c.Find(bson.M{"_id": bson.ObjectIdHex(team)}).One(&t)
	})
	return t.Accounts, err
}

func (r *team) GetVCSByID(team, id string) (types.VCS, error) {

	var t types.Team
	var err error
	r.Execute(func(c *mgo.Collection) {
		err = c.Find(bson.M{"name": team, "accounts.id": id}).Select(bson.M{"accounts.$": 1}).One(&t)
	})

	if len(t.Accounts) == 0 {
		return types.VCS{}, err
	}

	return t.Accounts[0], err
}

func (r *team) UpdateVCS(team string, vcs types.VCS) error {

	var err error
	r.Execute(func(c *mgo.Collection) {
		err = c.Update(bson.M{"name": team, "accounts.id": vcs.ID},
			bson.M{"$set": bson.M{
				"accounts.$.owner_type":   vcs.OwnerType,
				"accounts.$.token_expiry": vcs.TokenExpiry}})
	})
	return err
}

func (s *team) GetVCSByName(team, name, source string) (*types.VCS, error) {

	var err error
	var t types.Team
	s.Execute(func(c *mgo.Collection) {
		c.Find(bson.M{"name": team, "accounts.name": name, "accounts.source": source}).Select(bson.M{"accounts.$": 1}).One(&t)
	})

	if len(t.Accounts) == 0 {
		return nil, fmt.Errorf("Not found")
	}
	return &t.Accounts[0], err
}
