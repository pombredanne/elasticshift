// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

import (
	"testing"
	"time"

	"github.com/spf13/viper"

	"gopkg.in/mgo.v2"
	"gitlab.com/conspico/esh/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestTeamDatastore(t *testing.T) {
	suite.Run(t, new(TeamDatastoreTestSuite))
}

type TeamDatastoreTestSuite struct {
	suite.Suite
	config Config
	session *mgo.Session
	ds TeamDatastore

	//test specific
	id string
	vcs VCS
	team string
	t Team
}

func (suite *TeamDatastoreTestSuite) SetupTest() {

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

	suite.ds = NewTeamDatastore(ds)

	suite.team = "testteam"

}

func (suite *TeamDatastoreTestSuite) TearDownTest() {
	suite.session.Close()
}

func (suite *TeamDatastoreTestSuite) Test01Save() {

	team := Team{
		Name:    suite.team,
		Display: suite.team,
	}
	err := suite.ds.Save(&team)

	assert.Nil(suite.T(), err)
}

func (suite *TeamDatastoreTestSuite) Test02CheckExist() {

	exist, err := suite.ds.CheckExists(suite.team)

	assert.Nil(suite.T(), err)
	assert.True(suite.T(), exist)
}

func (suite *TeamDatastoreTestSuite) Test03GetTeam() {

	t, err := suite.ds.GetTeam(suite.team)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.team, t.Name)
}

func (suite *TeamDatastoreTestSuite) Test04SaveVCS() {

	v := &VCS{}
	v.ID = "vcsid"
	v.Name = "ghazninattarshah"
	v.OwnerType = "user"
	v.Type = "github"
	v.AccessCode = "accesscode"
	v.AccessToken = "accesstoken"
	v.AvatarURL = "avatar_url"
	v.RefreshToken = "refresh_token"
	v.TokenType = "Bearer"
	v.TokenExpiry = time.Now()

	err := suite.ds.SaveVCS(suite.team, v)
	assert.Nil(suite.T(), err)
}

func (suite *TeamDatastoreTestSuite) Test05GetVCS() {

	vcs, err := suite.ds.GetVCSByID(suite.team, "vcsid")
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), "ghazninattarshah", vcs.Name)

	suite.vcs = vcs
}

func (suite *TeamDatastoreTestSuite) Test06UpdateVCS() {

	suite.vcs.AccessCode = "updated_access_token"
	err := suite.ds.UpdateVCS(suite.team, suite.vcs)
	assert.Nil(suite.T(), err)
}

/*func TestTeam(t *testing.T) {

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
		t.Error(err)
		t.FailNow()
	}
	defer session.Close()

	ds := core.NewDatasource(config.DB.Name, session)

	teamDS := NewTeamDatastore(ds)

	team := Team{
		Name:    "test",
		Display: "test",
	}
	exist, err := teamDS.CheckExists(team.Name)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if exist {
		team, err = teamDS.GetTeam(team.Name)
	} else {
		err = teamDS.Save(&team)
	}

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	id, _ := util.NewUUID()
	v := &VCS{}
	v.ID = id
	v.Name = "ghazninattarshah"
	v.OwnerType = "user"
	v.Type = "github"
	v.AccessCode = "accesscode"
	v.AccessToken = "accesstoken"
	v.AvatarURL = "avatar_url"
	v.RefreshToken = "refresh_token"
	v.TokenType = "Bearer"
	v.TokenExpiry = time.Now()

	//v.ID = "92062b97e962460361426e193a5cdefb"

	err = teamDS.SaveVCS(team.Name, v)
	if err != nil {
		t.Log(err)
	}
	//var vnew esh.VCS
	vnew, err := teamDS.GetVCSByID(team.Name, v.ID)
	if err != nil {
		t.Log(err)
	}
	t.Log("List = ", vnew)

	vnew.AccessToken = "new_accesstoken"
	err = teamDS.UpdateVCS(team.Name, vnew)
	if err != nil {
		t.Log(err)
	}
}
*/