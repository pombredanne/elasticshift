package esh

type teamDatastore struct {
	ds Datastore
}

var (
	sqlCheckTeamExist = "SELECT 1 as 'exist' FROM TEAM WHERE name = ? LIMIT 1"

	sqlGetTeamID = "SELECT id FROM TEAM WHERE name = ? LIMIT 1"
)

func (r *teamDatastore) Save(team *Team) error {
	return r.ds.Create(team)
}

func (r *teamDatastore) CheckExists(name string) (bool, error) {

	var result struct {
		Exist int
	}

	err := r.ds.Read(sqlCheckTeamExist, &result, name)
	return result.Exist == 1, err
}

func (r *teamDatastore) GetTeamID(name string) (string, error) {

	var result struct {
		ID string
	}
	err := r.ds.Read(sqlGetTeamID, &result, name)
	return result.ID, err
}

// NewTeamDatastore ..
func NewTeamDatastore(ds Datastore) TeamDatastore {
	return &teamDatastore{ds: ds}
}
