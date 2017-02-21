package store

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type teamStore struct {
	store Store  // store
	cname string // collection name
}

// NewTeamStore related database operations
func NewTeamStore(s Store) TeamStore {
	return &teamStore{s, "team"}
}

// TeamStore related database operations
type TeamStore interface {
	//Team
	Save(team *Team) error
	CheckExists(name string) (bool, error)
	GetTeam(name string) (Team, error)

	// VCS Settings
	SaveVCS(team string, vcs *VCS) error
	UpdateVCS(team string, vcs VCS) error
	GetVCSByID(team, id string) (VCS, error)
}

func (r *teamStore) Save(team *Team) error {
	return r.store.Insert(r.cname, team)
}

func (r *teamStore) CheckExists(name string) (bool, error) {

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

func (r *teamStore) GetTeam(name string) (Team, error) {

	var err error
	var result Team
	r.store.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Find(bson.M{"name": name}).One(&result)
	})
	return result, err
}

func (r *teamStore) SaveVCS(team string, vcs *VCS) error {

	var err error
	r.store.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Update(
			bson.M{"name": team},
			bson.M{"$push": bson.M{"accounts": vcs}},
		)
	})
	return err
}

func (r *teamStore) GetVCSByID(team, id string) (VCS, error) {

	var t Team
	var err error
	r.store.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Find(bson.M{"name": team, "accounts._id": id}).Select(bson.M{"accounts.$": 1}).One(&t)
	})

	if len(t.Accounts) == 0 {
		return VCS{}, err
	}

	return t.Accounts[0], err
}

func (r *teamStore) UpdateVCS(team string, vcs VCS) error {

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
