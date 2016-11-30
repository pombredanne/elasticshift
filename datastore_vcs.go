package esh

type vcsDatastore struct {
	ds Datastore
}

var (
	sqlGetVCS = `SELECT id, 
							name, 
							type, 
							avatar_url,
							updated_dt
				     FROM VCS WHERE TEAM_ID = ?`

	sqlGetVCSByID = `SELECT id, 
			                team_id,
							name, 
							type, 
							avatar_url,
							access_token,
							refresh_token,
							token_expiry,
							token_type,
							owner_type
				     FROM VCS WHERE ID = ? LIMIT 1`

	sqlCheckVCSExist = "SELECT 1 as 'exist' FROM VCS WHERE TEAM_ID = ? AND VCS_ID = ? LIMIT 1"

	sqlGetByProvVCSID = `SELECT id, 
			                team_id,
							name, 
							type, 
							avatar_url,
							access_token,
							refresh_token,
							token_expiry,
							token_type,
							owner_type,
							vcs_id
				     FROM VCS WHERE TEAM_ID = ? AND VCS_ID = ? LIMIT 1`
)

func (r *vcsDatastore) Save(v *VCS) error {
	return r.ds.Create(v)
}

func (r *vcsDatastore) GetVCS(teamID string) ([]VCS, error) {

	var result []VCS
	err := r.ds.Read(sqlGetVCS, &result, teamID)
	return result, err
}

func (r *vcsDatastore) GetByID(id string) (VCS, error) {

	var result VCS
	err := r.ds.Read(sqlGetVCSByID, &result, id)
	return result, err
}

func (r *vcsDatastore) Update(old *VCS, updated VCS) error {
	return r.ds.Update(&old, updated)
}

func (r *vcsDatastore) CheckIfExists(vcsID, teamID string) (bool, error) {

	var result struct {
		Exist int
	}

	err := r.ds.Read(sqlCheckVCSExist, &result, vcsID, teamID)
	return result.Exist == 1, err
}

func (r *vcsDatastore) GetByProviderVCSID(teamID, vcsID string) (VCS, error) {

	var result VCS
	err := r.ds.Read(sqlGetByProvVCSID, &result, teamID, vcsID)
	return result, err
}

// NewVCSDatastore ..
func NewVCSDatastore(ds Datastore) VCSDatastore {
	return &vcsDatastore{ds: ds}
}
