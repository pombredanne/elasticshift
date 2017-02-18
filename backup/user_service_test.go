// Package esh
// Author Ghazni Nattarshah
// Date: Jan 3, 2017
package esh

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/conspico/esh/core"
	"gitlab.com/conspico/esh/core/util"
	"gopkg.in/mgo.v2"
	"testing"
)

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

type UserServiceTestSuite struct {
	suite.Suite
	config  Config
	session *mgo.Session
	svc     UserService
	appCtx  AppContext

	//test specific
	name     string
	email    string
	password string
}

func (suite *UserServiceTestSuite) SetupTest() {

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

	signer, err := util.LoadKey(config.Key.Signer)
	assert.Nil(suite.T(), err)

	appCtx := AppContext{}
	appCtx.Signer = signer
	appCtx.Datasource = core.NewDatasource(config.DB.Name, session)
	appCtx.UserDatastore = NewUserDatastore(appCtx.Datasource)
	appCtx.TeamDatastore = NewTeamDatastore(appCtx.Datasource)
	appCtx.Config = config
	appCtx.Logger = logrus.New()

	suite.appCtx = appCtx
	suite.svc = NewUserService(appCtx)

	if suite.name == "" {

		suite.name, _ = util.NewUUID()
		suite.email = suite.name + "@email.com"
	}
	suite.password = "T3$tu$3r01"
}

func (suite *UserServiceTestSuite) TearDownTest() {
	suite.session.Close()
}

//Create(r signupRequest) (string, error)
//SignIn(r signInRequest) (string, error)
//SignOut() (bool, error)
//Verify(code string) (bool, error)
func (suite *UserServiceTestSuite) Test01Create() {

	req := signupRequest{}
	req.Team = suite.name
	req.Fullname = suite.name
	req.Email = suite.email
	req.Password = suite.password
	created, err := suite.svc.Create(req)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), created)
}

func (suite *UserServiceTestSuite) Test02CreateUserAlreadyExist() {

	req := signupRequest{}
	req.Team = suite.name
	req.Fullname = suite.name
	req.Email = suite.email
	req.Password = suite.password
	created, err := suite.svc.Create(req)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "", created)
}

func (suite *UserServiceTestSuite) Test03CreateUserNoPassword() {

	req := signupRequest{}
	req.Team = suite.name
	req.Fullname = suite.name
	req.Email = suite.email
	req.Password = ``
	created, err := suite.svc.Create(req)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "", created)
}

func (suite *UserServiceTestSuite) Test04Sign() {

	req := signInRequest{}
	req.Team = suite.name
	req.Email = suite.email
	req.Password = suite.password

	tok, err := suite.svc.SignIn(req)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), tok)
}

func (suite *UserServiceTestSuite) Test05SignInvalidUsername() {

	req := signInRequest{}
	req.Team = suite.name
	req.Email, _ = util.NewUUID()
	req.Password = suite.password

	tok, err := suite.svc.SignIn(req)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "", tok)
}

func (suite *UserServiceTestSuite) Test06SignInvalidPassword() {

	req := signInRequest{}
	req.Team = suite.name
	req.Email = suite.email
	req.Password = "password"

	tok, err := suite.svc.SignIn(req)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), errInvalidEmailOrPassword.Error(), tok)
}

func (suite *UserServiceTestSuite) Test07Signout() {

	success, err := suite.svc.SignOut()

	assert.Nil(suite.T(), err)
	assert.True(suite.T(), success)
}
