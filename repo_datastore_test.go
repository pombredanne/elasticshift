// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

import (
	"testing"

	"github.com/spf13/viper"

	"gopkg.in/mgo.v2"
)

func TestDatastore(t *testing.T) {

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
	}
	defer session.Close()

	ds := NewDatasource(config.DB.Name, session)

	repoDS := NewRepoDatastore(ds)

	r, err := repoDS.GetReposByVCSID("conspico", "7168293")
	t.Log(r)
}
