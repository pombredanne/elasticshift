// Package esh ...
// Author: Ghazni Nattarshah
// Date: Dec 30, 2016
package core

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	recordNotFound = "record not found"
)

// Datastore ..
// Abstract Datastore to interact with DB.
type Datastore interface {
	Execute(cname string, handleFunc func(c *mgo.Collection))
	Insert(cname string, model interface{}) error
	Upsert(cname string, selector interface{}, model interface{}) (*mgo.ChangeInfo, error)
	FindAll(cname string, query interface{}, model interface{}) error
	FindOne(cname string, query interface{}, model interface{}) error
	Exist(cname string, selector interface{}) (bool, error)
	Remove(cname string, id bson.ObjectId) error
	RemoveMultiple(cname string, ids []bson.ObjectId) error
}

// Datastore ..
// A base datasource that performs actualy sql interactions.
type datastore struct {
	session  *mgo.Session
	database string
}

// Executes a given func with a active session against the database.
func (ds *datastore) Execute(cname string, handle func(c *mgo.Collection)) {

	s := ds.session.Copy()
	defer s.Close()

	handle(s.DB(ds.database).C(cname))
	return
}

// Checks whether the given document exist in a collection
func (ds *datastore) Exist(cname string, selector interface{}) (bool, error) {

	var count int
	var err error
	ds.Execute(cname, func(c *mgo.Collection) {
		count, err = c.Find(selector).Count()
	})

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Performs and insert operation for a model on a collection
func (ds *datastore) Insert(cname string, model interface{}) error {

	var err error
	ds.Execute(cname, func(c *mgo.Collection) {
		err = c.Insert(model)
	})
	return err
}

// Performs Upsert for a model on a collection, based on a selector
func (ds *datastore) Upsert(cname string, selector interface{}, model interface{}) (*mgo.ChangeInfo, error) {

	var info *mgo.ChangeInfo
	var err error
	ds.Execute(cname, func(c *mgo.Collection) {
		info, err = c.Upsert(selector, model)
	})
	return info, err
}

// Find all the document matches the query on a collection.
func (ds *datastore) FindAll(cname string, query interface{}, model interface{}) error {

	var err error
	ds.Execute(cname, func(c *mgo.Collection) {
		err = c.Find(query).All(model)
	})
	return err
}

// Find one document matches the query on a collection
func (ds *datastore) FindOne(cname string, query interface{}, model interface{}) error {

	var err error
	ds.Execute(cname, func(c *mgo.Collection) {
		err = c.Find(query).One(model)
	})
	return err
}

// Remove a document based on id
func (ds *datastore) Remove(cname string, id bson.ObjectId) error {

	var err error
	ds.Execute(cname, func(c *mgo.Collection) {
		err = c.RemoveId(id)
	})
	return err
}

// Removes multiple document based on gived ids
func (ds *datastore) RemoveMultiple(cname string, ids []bson.ObjectId) error {

	var err error
	ds.Execute(cname, func(c *mgo.Collection) {
		err = c.Remove(bson.M{"_id": bson.M{"$in": ids}})
	})
	return err
}

// NewDatasource ..
// Create a new base datasource
func NewDatasource(dbname string, session *mgo.Session) Datastore {
	return &datastore{database: dbname, session: session}
}
