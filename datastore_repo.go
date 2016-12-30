// Package esh ...
// Author: Ghazni Nattarshah
// Date: Oct 22, 2016
package esh

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type repoDatastore struct {
	ds    Datastore
	cname string
}

func (r *repoDatastore) Save(repo *Repo) error {
	return r.ds.Insert(r.cname, repo)
}

func (r *repoDatastore) GetReposByVCSID(team, vcsID string) ([]Repo, error) {

	result := []Repo{}
	r.ds.FindAll(r.cname, bson.M{"team": team, "vcs_id": vcsID}, &result)
	return result, nil
}

func (r *repoDatastore) Update(repo Repo) error {
	var err error
	r.ds.Execute(r.cname, func(c *mgo.Collection) {
		err = c.Update(bson.M{"_id": repo.ID},
			bson.M{"$set": bson.M{"name": repo.Name,
				"private":        repo.Private,
				"link":           repo.Link,
				"description":    repo.Description,
				"fork":           repo.Fork,
				"default_branch": repo.DefaultBranch,
				"language":       repo.Language}})
	})
	return err
}

func (r *repoDatastore) Delete(repo Repo) error {
	return r.ds.Remove(r.cname, repo.ID)
}

func (r *repoDatastore) DeleteIds(ids []bson.ObjectId) error {
	return r.ds.RemoveMultiple(r.cname, ids)
}

func (r *repoDatastore) GetRepos(team string) ([]Repo, error) {

	var result []Repo
	err := r.ds.FindAll(r.cname, bson.M{"team": team}, &result)
	return result, err
}

// NewRepoDatastore ..
func NewRepoDatastore(ds Datastore) RepoDatastore {
	return &repoDatastore{ds: ds, cname: "repos"}
}
