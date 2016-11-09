package esh

type repoService struct {
	repoDS RepoDatastore
	config Config
}

// NewRepoService ..
func NewRepoService(r RepoDatastore, conf Config) RepoService {

	return &repoService{
		repoDS: r,
		config: conf,
	}
}

func (s *repoService) GetRepos(teamID string) (GetRepoResponse, error) {

	result, err := s.repoDS.GetRepos(teamID)
	return GetRepoResponse{Result: result}, err
}

func (s *repoService) GetReposByVCSID(teamID, vcsID string) (GetRepoResponse, error) {

	result, err := s.repoDS.GetReposByVCSID(teamID, vcsID)
	return GetRepoResponse{Result: result}, err
}
