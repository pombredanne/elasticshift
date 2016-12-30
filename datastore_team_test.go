// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

import (
	"testing"
	"time"

	"gitlab.com/conspico/esh"
	"gitlab.com/conspico/esh/core/util"

	"github.com/spf13/viper"

	"gopkg.in/mgo.v2"
)

func TestLoadSysconf(t *testing.T) {

	vip := viper.New()
	vip.SetConfigType("yml")
	vip.SetConfigFile("esh.yml")
	vip.ReadInConfig()

	config := esh.Config{}
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

	ds := esh.NewDatasource(config.DB.Name, session)

	teamDS := esh.NewTeamDatastore(ds)

	team := esh.Team{
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
	v := &esh.VCS{}
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
