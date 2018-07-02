/*
Copyright 2018 The Elasticshift Authors.
*/
package plugin

import (
	"context"
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errIDCantBeEmpty    = errors.New("Plugin ID cannot be empty")
	errNameCantBeEmpty  = errors.New("Plugin Name cannot be empty")
	errTeamOrNameIsMust = errors.New("Plugin name or team name is must")
)

// Resolver ...
type Resolver interface {
	FetchPlugin(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	store  store.Plugin
	logger logrus.Logger
	Ctx    context.Context
}

// NewResolver ...
func NewResolver(ctx context.Context, logger logrus.Logger, s store.Shift) (Resolver, error) {

	r := &resolver{
		store:  s.Plugin,
		logger: logger,
		Ctx:    ctx,
	}
	return r, nil
}

func (r *resolver) FetchPlugin(params graphql.ResolveParams) (interface{}, error) {

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
	var result []types.Plugin
	r.store.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	var res types.PluginList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}
