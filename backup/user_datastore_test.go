// Package esh
// Author Ghazni Nattarshah
// Date: 1/2/17
package esh

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/conspico/esh/core"
	"gopkg.in/mgo.v2"
)

func TestUserDatastore(t *testing.T) {
	suite.Run(t, new(UserDatastoreTestSuite))
}

type UserDatastoreTestSuite struct {
	suite.Suite
	config  Config
	session *mgo.Session
	ds      UserDatastore

	//test specific
	id    string
	team  string
	email string
}

func (suite *UserDatastoreTestSuite) SetupTest() {

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

	suite.ds = NewUserDatastore(ds)

	suite.team = "testteam"
	suite.email = "test.user@email.com"
}

func (suite *UserDatastoreTestSuite) TearDownTest() {
	suite.session.Close()
}

func (suite *UserDatastoreTestSuite) Test01Save() {

	user := User{}
	user.Fullname = "Test User"
	user.Active = true
	user.BadAttempt = 0
	user.Email = suite.email
	user.EmailVefified = true
	user.Locked = false
	user.Password = "Password1!"
	user.Team = suite.team

	suite.ds.Save(&user)
	err := suite.ds.Save(&user)

	assert.Nil(suite.T(), err)
}

func (suite *UserDatastoreTestSuite) Test02CheckExist() {

	exist, err := suite.ds.CheckExists(suite.email, suite.team)

	assert.Nil(suite.T(), err)
	assert.True(suite.T(), exist)
}

func (suite *UserDatastoreTestSuite) Test03GetUser() {

	u, err := suite.ds.GetUser(suite.email, suite.team)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), suite.email, u.Email)
}
