/*
Copyright 2018 The Elasticshift Authors.
*/
package app

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
	errIDCantBeEmpty    = errors.New("App ID cannot be empty")
	errNameCantBeEmpty  = errors.New("App Name cannot be empty")
	errTeamOrNameIsMust = errors.New("App name or team name is must")
)

type resolver struct {
	store  Store
	logger logrus.Logger
	Ctx    context.Context
}

func (r *resolver) FetchApp(params graphql.ResolveParams) (interface{}, error) {

	name := params.Args["name"].(string)
	team := params.Args["team"].(string)

	if name == "" && team == "" {
		return nil, errTeamOrNameIsMust
	}

	q := bson.M{}

	if team != "" {
		q["team"] = team
	}

	if name != "" {
		q["name"] = name
	}

	var err error
	var result []types.App
	r.store.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	var res types.AppList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}
