// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

import (
	"testing"

	"github.com/spf13/viper"

	"gopkg.in/mgo.v2"
	"github.com/stretchr/testify/suite"
	"gitlab.com/conspico/esh/core"
	"gitlab.com/conspico/esh/core/util"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestRepoDatastore(t *testing.T) {
	suite.Run(t, new(RepoDatastoreTestSuite))
}

type RepoDatastoreTestSuite struct {
	suite.Suite
	config Config
	session *mgo.Session
	repoDS RepoDatastore

	//test specific
	id string
	repo Repo
	team string
	vcsid string
}

func (suite *RepoDatastoreTestSuite) SetupTest() {

	vip := viper.New()
	vip.SetConfigType("yml")
	vip.SetConfigFile("esh.yml")
	vip.ReadInConfig()

	config := Config{}
	vip.Unmarshal(&config)

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{config.DB.Server},
		Username: config.DB.Username,
		Password: config.DB.Password,
		Database: config.DB.Name,
	})

	if err != nil {
		suite.T().Log(err)
		suite.T().FailNow()
	}

	suite.session = session
	ds := core.NewDatasource(config.DB.Name, session)

	suite.repoDS = NewRepoDatastore(ds)

	suite.team = "testteam"
	suite.vcsid = "testvcs"

}

func (suite *RepoDatastoreTestSuite) TearDownTest() {
	suite.session.Close()
}

func (suite *RepoDatastoreTestSuite) Test01Save() {


	repo := Repo{}
	suite.id, _ = util.NewUUID()
	repo.Team = suite.team
	repo.DefaultBranch = "develop"
	repo.Description = "test project"
	repo.Fork = true
	repo.Language = "Java"
	repo.Link = "http://test.project.com"
	repo.Name = "testproject"
	repo.RepoID = "12345"
	repo.VcsID = suite.vcsid

	err := suite.repoDS.Save(&repo)
	assert.Nil(suite.T(), err)
}

func (suite *RepoDatastoreTestSuite) Test02GetRepos() {

	repos, err := suite.repoDS.GetRepos(suite.team)

	assert.Nil(suite.T(), err)
	assert.True(suite.T(), len(repos) > 0)

	suite.repo = repos[0]
}

func (suite *RepoDatastoreTestSuite) Test03Update() {

	suite.repo.Description = "Updated description"

	err := suite.repoDS.Update(suite.repo)
	assert.Nil(suite.T(), err)
}

func (suite *RepoDatastoreTestSuite) Test04GetReposByVCSID() {

	repos, err := suite.repoDS.GetReposByVCSID(suite.team, suite.vcsid)

	assert.Nil(suite.T(), err)
	assert.True(suite.T(), len(repos) > 0)
}

func (suite *RepoDatastoreTestSuite) Test05Delete() {

	err := suite.repoDS.Delete(suite.repo)
	assert.Nil(suite.T(), err)
}

func (suite *RepoDatastoreTestSuite) Test06DeleteIds() {

	err := suite.repoDS.Save(&suite.repo)
	assert.Nil(suite.T(), err)

	var repos []Repo
	repos, err = suite.repoDS.GetRepos(suite.team)

	var ids []bson.ObjectId
	for _, r := range repos {
		ids = append(ids, r.ID)
	}

	err = suite.repoDS.DeleteIds(ids)
	assert.Nil(suite.T(), err)
}
