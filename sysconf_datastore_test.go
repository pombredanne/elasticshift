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
	"gopkg.in/mgo.v2/bson"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestSysconfDatastore(t *testing.T) {
	suite.Run(t, new(SysconfDatastoreTestSuite))
}

type SysconfDatastoreTestSuite struct {
	suite.Suite
	config Config
	session *mgo.Session
	ds SysconfDatastore

	//test specific
	id bson.ObjectId
	repo VCSSysConf
	team string
}

func (suite *SysconfDatastoreTestSuite) SetupTest() {

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

	suite.ds = NewSysconfDatastore(ds)

	suite.team = "testteam"

}

func (suite *SysconfDatastoreTestSuite) TearDownTest() {
	suite.session.Close()
}

func (suite *SysconfDatastoreTestSuite) Test01Save() {


	sc := &VCSSysConf{}
	sc.Name = "ghub"
	sc.Type = "vcs"
	sc.Secret = "secret"
	sc.Key = "key"
	sc.CallbackURL = "callback_url"
	sc.HookURL = "hook_url"

	err := suite.ds.SaveVCS(sc)
	assert.Nil(suite.T(), err)
}

func (suite *SysconfDatastoreTestSuite) Test02GetVCSTypes() {

	scf, err := suite.ds.GetVCSTypes()

	assert.Nil(suite.T(), err)
	assert.True(suite.T(), len(scf) > 0)

	for _, v := range scf {
		if strings.EqualFold("ghub", v.Name) {
			suite.id = v.ID
			break;
		}
	}
}

func (suite *SysconfDatastoreTestSuite) Test03Delete() {

	err := suite.ds.Delete(suite.id)
	assert.Nil(suite.T(), err)
}