package esh

type userDatastore struct {
	ds Datastore
}

var (
	sqlCheckUserExist = "SELECT 1 as 'exist' FROM USER WHERE email = ? AND TEAM_ID = ? LIMIT 1"

	sqlGetUser = `SELECT id, 
							team_id, 
							fullname, 
							username,
							email,
							password,
							locked,
							active,
							bad_attempt
				     FROM USER WHERE email = ? AND TEAM_ID = ? LIMIT 1`
)

func (r *userDatastore) Save(user *User) error {
	return r.ds.Create(user)
}

func (r *userDatastore) CheckExists(email, teamID string) (bool, error) {

	var result struct {
		Exist int
	}

	err := r.ds.Read(sqlCheckUserExist, &result, email, teamID)
	return result.Exist == 1, err
}

func (r *userDatastore) GetUser(email, teamID string) (User, error) {

	var result User
	err := r.ds.Read(sqlGetUser, &result, email, teamID)
	return result, err
}

// NewUserDatastore ..
func NewUserDatastore(ds Datastore) UserDatastore {
	return &userDatastore{ds: ds}
}
