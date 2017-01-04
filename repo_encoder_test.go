// Package esh
// Author Ghazni Nattarshah
// Date: Jan 4, 2017
package esh

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"context"
	"net/http/httptest"
)



func TestRepoEncoder(t *testing.T) {
	suite.Run(t, new(RepoEncoderTestSuite))
}

type RepoEncoderTestSuite struct {
	suite.Suite
	ctx context.Context
}

func (suite *RepoEncoderTestSuite) TestGetRepoRequest() {

	repo := Repo{}
	repo.Team = "testteam"
	repo.DefaultBranch = "develop"
	repo.Description = "test project"
	repo.Fork = true
	repo.Language = "Java"
	repo.Link = "http://test.project.com"
	repo.Name = "testproject"
	repo.RepoID = "12345"


	res := GetRepoResponse{}
	res.Result = []Repo{repo}

	err := encodeGetRepoResponse(context.Background(), httptest.NewRecorder(), res)
	assert.Nil(suite.T(), err)
}