package store

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type repositoryStore struct {
	store Store  // store
	cname string // collection name
}

// NewRepositoryStore related database operations
func NewRepositoryStore(s Store) RepositoryStore {
	return &repositoryStore{s, "repository"}
}

// RepositoryStore related database operations
type RepositoryStore interface {
	Save(repo *Repository) error
	Update(repo Repository) error
	Delete(repo Repository) error
	DeleteIds(ids []interface{}) error
	GetRepos(teamID string) ([]Repository, error)
	GetReposByVCSID(team, vcsID string) ([]Repository, error)
}

func (r *repositoryStore) Save(repo *Repository) error {
	return r.store.Insert(r.cname, repo)
}

func (r *repositoryStore) GetReposByVCSID(team, vcsID string) ([]Repository, error) {

	result := []Repository{}
	r.store.FindAll(r.cname, bson.M{"team": team, "vcs_id": vcsID}, &result)
	return result, nil
}

func (r *repositoryStore) Update(repo Repository) error {
	var err error
	r.store.Execute(r.cname, func(c *mgo.Collection) {
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

func (r *repositoryStore) Delete(repo Repository) error {
	return r.store.Remove(r.cname, repo.ID)
}

func (r *repositoryStore) DeleteIds(ids []interface{}) error {
	return r.store.RemoveMultiple(r.cname, ids)
}

func (r *repositoryStore) GetRepos(team string) ([]Repository, error) {

	var result []Repository
	err := r.store.FindAll(r.cname, bson.M{"team": team}, &result)
	return result, err
}
