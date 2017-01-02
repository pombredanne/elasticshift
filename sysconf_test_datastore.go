// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package esh

import (
	"testing"

	"github.com/spf13/viper"

	"gopkg.in/mgo.v2"
	"gitlab.com/conspico/esh/core"
)

func TestLoadSysconf(t *testing.T) {

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

	ds := core.NewDatasource(config.DB.Name, session)

	syscDS := NewSysconfDatastore(ds)

	/*sc := &esh.Sysconf{}
	sc.ID = bson.NewObjectIdWithTime(time.Now())
	sc.Name = "ghub"
	sc.Type = "vcs"

	vd := &esh.VCSData{}
	vd.Secret = "secret"
	vd.Key = "key"
	vd.CallbackURL = "callback_url"
	vd.HookURL = "hook_url"
	data, _ := json.Marshal(vd)
	sc.Data = data
	syscDS.Save(sc)
	*/

	if err != nil {
		t.Error(err)
	}
	scf, err := syscDS.GetVCSTypes()
	if err != nil {
		t.Error(err)
	}
	for _, v := range scf {
		t.Log(v.Name)
	}
	//t.Log(scf)
}
