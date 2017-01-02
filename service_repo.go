// Package esh ...
// Author: Ghazni Nattarshah
// Date: Oct 29, 2016
package esh

type repoService struct {
	repoDS RepoDatastore
	config Config
}

// NewRepoService ..
func NewRepoService(appCtx AppContext) RepoService {

	return &repoService{
		repoDS: appCtx.RepoDatastore,
		config: appCtx.Config,
	}
}

func (s *repoService) GetRepos(team string) (GetRepoResponse, error) {

	result, err := s.repoDS.GetRepos(team)
	return GetRepoResponse{Result: result}, err
}

func (s *repoService) GetReposByVCSID(team, id string) (GetRepoResponse, error) {

	result, err := s.repoDS.GetReposByVCSID(team, id)
	return GetRepoResponse{Result: result}, err
}
