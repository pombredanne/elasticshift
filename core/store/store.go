/*
Copyright 2017 The Elasticshift Authors.
*/
package store

import (
	"encoding/json"

	"gitlab.com/conspico/elasticshift/core/utils"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	recordNotFound = "record not found"
)

// Store ..
// Abstract database interactions.
type Store interface {
	Execute(cname string, handleFunc func(c *mgo.Collection))
	Insert(cname string, model interface{}) error
	Upsert(cname string, selector interface{}, model interface{}) (*mgo.ChangeInfo, error)
	FindAll(cname string, query interface{}, model interface{}) error
	FindOne(cname string, query interface{}, model interface{}) error
	Exist(cname string, selector interface{}) (bool, error)
	Remove(cname string, id interface{}) error
	RemoveMultiple(cname string, ids []interface{}) error
	GetSession() *mgo.Session
}

// New ..
// Create a new base datasource
func New(dbname string, session *mgo.Session) Store {
	return &store{database: dbname, session: session}
}

// Store ..
// A base datasource that performs actualy sql interactions.
type store struct {
	session  *mgo.Session
	database string
}

func (s *store) GetSession() *mgo.Session {
	return s.session
}

// Execute given func with a active session against the database
func (s *store) Execute(cname string, handle func(c *mgo.Collection)) {

	ses := s.session.Copy()
	defer ses.Close()

	handle(ses.DB(s.database).C(cname))
	return
}

// Checks whether the given document exist in a collection
func (s *store) Exist(cname string, selector interface{}) (bool, error) {

	var count int
	var err error
	s.Execute(cname, func(c *mgo.Collection) {
		count, err = c.Find(selector).Count()
	})

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Insert operation for a model on a collection
func (s *store) Insert(cname string, model interface{}) error {

	var err error
	s.Execute(cname, func(c *mgo.Collection) {
		err = c.Insert(model)
	})
	return err
}

// Upsert for a model on a collection, based on a selector
func (s *store) Upsert(cname string, selector interface{}, model interface{}) (*mgo.ChangeInfo, error) {

	var info *mgo.ChangeInfo
	var err error
	s.Execute(cname, func(c *mgo.Collection) {
		info, err = c.Upsert(selector, model)
	})
	return info, err
}

// FindAll the document matches the query on a collection.
func (s *store) FindAll(cname string, query interface{}, model interface{}) error {

	var err error
	s.Execute(cname, func(c *mgo.Collection) {
		err = c.Find(query).All(model)
	})
	return err
}

// FindOne document matches the query on a collection
func (s *store) FindOne(cname string, query interface{}, model interface{}) error {

	var err error
	s.Execute(cname, func(c *mgo.Collection) {
		err = c.Find(query).One(model)
	})
	return err
}

// Remove a document based on id
func (s *store) Remove(cname string, id interface{}) error {

	var err error
	s.Execute(cname, func(c *mgo.Collection) {
		err = c.RemoveId(id)
	})
	return err
}

// RemoveMultiple document based on gived ids
func (s *store) RemoveMultiple(cname string, ids []interface{}) error {

	var err error
	s.Execute(cname, func(c *mgo.Collection) {
		err = c.Remove(bson.M{"_id": bson.M{"$in": ids}})
	})
	return err
}

func encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func decode(data []byte, out interface{}) error {
	return json.Unmarshal(data, out)
}

// NewID ..
// Creates a new UUID and returns string
func NewID() string {
	return utils.NewUUID()
}
