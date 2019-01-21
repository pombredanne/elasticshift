/*
Copyright 2019 The Elasticshift Authors.
*/
package testhelper

import (
	"io/ioutil"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/dbtest"
)

var dbserv dbtest.DBServer

// ConnectMongo ...
func ConnectMongo() *mgo.Session {

	tempDir, _ := ioutil.TempDir("", "testing")
	dbserv.SetPath(tempDir)

	return dbserv.Session()
}

// Close ...
func Close(session *mgo.Session) {
	session.Close()
}

// Stop ...
func Stop(session *mgo.Session) {
	Close(session)
	dbserv.Stop()
}
