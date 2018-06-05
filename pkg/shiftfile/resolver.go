/*
Copyright 2018 The Elasticshift Authors.
*/
package shiftfile

import (
	"context"
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errIDCantBeEmpty    = errors.New("Plugin ID cannot be empty")
	errNameCantBeEmpty  = errors.New("Plugin Name cannot be empty")
	errTeamOrNameIsMust = errors.New("Plugin name or team name is must")
)

type resolver struct {
	store  Store
	logger logrus.Logger
	Ctx    context.Context
}

func (r *resolver) FetchShiftfile(params graphql.ResolveParams) (interface{}, error) {

	name := params.Args["name"].(string)
	if name == "" {
		return nil, errNameCantBeEmpty
	}

	q := bson.M{}

	if name != "" {
		q["name"] = name
	}

	var err error
	var result []types.Shiftfile
	r.store.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	var res types.ShiftfileList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}
