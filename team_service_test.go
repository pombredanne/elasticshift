// Package esh
// Author Ghazni Nattarshah
// Date: Jan 3, 2017
package esh

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2"
	"github.com/spf13/viper"
	"gitlab.com/conspico/esh/core"
	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/esh/core/util"
)

func TestTeamService(t *testing.T) {
	suite.Run(t, new(TeamServiceTestSuite))
}

type TeamServiceTestSuite struct {
	suite.Suite
	config Config
	session *mgo.Session
	svc TeamService
	appCtx AppContext

	//test specific
	team string
}

func (suite *TeamServiceTestSuite) SetupTest() {

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
	appCtx.TeamDatastore = NewTeamDatastore(appCtx.Datasource)
	appCtx.Config = config
	appCtx.Logger = logrus.New()

	suite.appCtx = appCtx
	suite.svc = NewTeamService(appCtx)

	if suite.team == "" {
		suite.team, _ = util.NewUUID()
	}
}

func (suite *TeamServiceTestSuite) TearDownTest() {

	suite.session.Close()
}

func (suite *TeamServiceTestSuite) Test01Create() {

	created, err := suite.svc.Create(suite.team)

	assert.Nil(suite.T(), err)
	assert.True(suite.T(), created)
}

func (suite *TeamServiceTestSuite) Test02CreateTeamAlreadyExist() {

	created, err := suite.svc.Create(suite.team)

	assert.Equal(suite.T(), errTeamAlreadyExists, err)
	assert.False(suite.T(), created)
}
