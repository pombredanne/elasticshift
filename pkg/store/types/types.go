/*
Copyright 2017 The Elasticshift Authors.
*/
package types

import mgo "gopkg.in/mgo.v2"

type Database struct {
	Session *mgo.Session
	Name    string
}
