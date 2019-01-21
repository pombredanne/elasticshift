package store

import (
	"testing"
)

func TestFetchBuild(t *testing.T) {

	//shift := testInitShiftStore()

	// b := types.Build{}
	// b.ID = bson.NewObjectId()
	// b.RepositoryID = "repositoryid"
	// b.ContainerEngineID = def.ContainerEngineID
	// b.VcsID = repo.VcsID
	// b.TriggeredBy = "Anonymous" //TODO fill in with logged-in user
	// b.Team = repo.Team
	// b.Branch = branch
	// b.StorageID = def.StorageID
	// b.CloneURL = repo.CloneURL
	// b.Language = repo.Language
	// b.Source = repo.Source
	// sb := types.SubBuild{
	// 	ID:     "1",
	// 	Graph:  defaultGraph,
	// 	Status: status,
	// }
	// b.SubBuilds = []types.SubBuild{sb}

	// Build file path - (for NFS)
	// <cache>/team-id/vcs-id/repository-id/branch-name/build-id/log
	// <cache>/team-id/vcs-id/repository-id/branch-name/build-id/reports
	// <cache>/team-id/vcs-id/repository-id/branch-name/build-id/archive.zip
	// cache must be mounted as /elasticshift to containers
	// b.StoragePath = filepath.Join(repo.Team, repo.Identifier, repo.Name, branch)
}
