/*
Copyright 2018 The Elasticshift Authors.
*/
package store

import (
	"testing"

	"github.com/elasticshift/elasticshift/pkg/testhelper"
	mgo "gopkg.in/mgo.v2"
)

func TestInitDatabase(t *testing.T) {
	testInitShiftStore()
}

func testInitShiftStore() (Shift, *mgo.Session) {
	session := testhelper.ConnectMongo()
	db := NewDatabase("testdb", session)
	return db.InitShiftStore(), session
}
