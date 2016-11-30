package esh

type repoDatastore struct {
	ds Datastore
}

var (
	sqlGetReposByVCSID = `SELECT id,
							team_id,
							vcs_id,
							repo_id,
							name,
							private,
							link,
							description,
							fork,
							default_branch,
							language
						FROM REPO WHERE team_id = ? and vcs_id = ?`
	sqlGetRepos = `SELECT id,
							team_id,
							vcs_id,
							repo_id,
							name,
							private,
							link,
							description,
							fork,
							default_branch,
							language
					FROM REPO WHERE team_id = ?`
)

func (r *repoDatastore) Save(repo *Repo) error {
	return r.ds.Create(repo)
}

func (r *repoDatastore) GetReposByVCSID(teamID, vcsID string) ([]Repo, error) {

	result := []Repo{}
	err := r.ds.Read(sqlGetReposByVCSID, &result, teamID, vcsID)
	return result, err
}

func (r *repoDatastore) Update(old Repo, repo Repo) error {
	return r.ds.Update(&old, repo)
}

func (r *repoDatastore) Delete(repo Repo) error {
	return r.ds.Delete(&repo)
}

func (r *repoDatastore) DeleteIds(ids []string) error {
	return r.ds.DeleteMultiple(Repo{}, ids)
}

func (r *repoDatastore) GetRepos(teamID string) ([]Repo, error) {

	var result []Repo
	err := r.ds.Read(sqlGetRepos, &result, teamID)
	return result, err
}

// NewRepoDatastore ..
func NewRepoDatastore(ds Datastore) RepoDatastore {
	return &repoDatastore{ds: ds}
}
