// Package esh
// Author Ghazni Nattarshah
// Date: 1/3/17
package esh

import (
	"testing"

	"github.com/spf13/viper"

	"gopkg.in/mgo.v2"
	"github.com/stretchr/testify/suite"
	"gitlab.com/conspico/esh/core"
	"github.com/stretchr/testify/assert"
	"gitlab.com/conspico/esh/core/util"
)

func TestRepoService(t *testing.T) {
	suite.Run(t, new(RepoServiceTestSuite))
}

type RepoServiceTestSuite struct {
	suite.Suite
	config Config
	session *mgo.Session
	svc RepoService
	appCtx AppContext

	//test specific
	id string
	repo Repo
	team string
	vcsid string
}

func (suite *RepoServiceTestSuite) SetupTest() {

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

	appCtx := AppContext{}
	appCtx.Datasource = core.NewDatasource(config.DB.Name, session)
	appCtx.RepoDatastore = NewRepoDatastore(appCtx.Datasource)
	appCtx.Config = config
	suite.appCtx = appCtx

	suite.svc = NewRepoService(appCtx)

	suite.team = "testteam"
	suite.vcsid = "testvcs"

	//Create a repo
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

	err = appCtx.RepoDatastore.Save(&repo)
	assert.Nil(suite.T(), err)

}

func (suite *RepoServiceTestSuite) TearDownTest() {

	suite.appCtx.RepoDatastore.Delete(suite.repo)

	suite.session.Close()
}

func (suite *RepoServiceTestSuite) Test01GetRepos() {

	repos, err := suite.svc.GetRepos(suite.team)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), len(repos.Result) > 0)
}

func (suite *RepoServiceTestSuite) Test01GetReposBYVCSID() {

	repos, err := suite.svc.GetReposByVCSID(suite.team, suite.vcsid)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), len(repos.Result) > 0)

	suite.repo = repos.Result[0]

}