/*
Copyright 2018 The Elasticshift Authors.
*/
package container

import (
	"context"
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errIDCantBeEmpty     = errors.New("Container ID cannot be empty")
	errTeamCannotBeEmpty = errors.New("Team must be provided")
)

type resolver struct {
	store  store.Container
	logger logrus.Logger
	Ctx    context.Context
}

func (r *resolver) FetchContainer(params graphql.ResolveParams) (interface{}, error) {

	team := params.Args["team"].(string)
	if team == "" {
		return nil, errTeamCannotBeEmpty
	}

	q := bson.M{"team": team}

	id := params.Args["id"].(string)
	if id != "" {
		q["id"] = id
	}

	buildID := params.Args["build_id"].(string)
	if buildID != "" {
		q["build_id"] = buildID
	}

	status := params.Args["status"].(int)
	if status > 0 {
		q["status"] = status
	}

	var err error
	var result []types.Container
	r.store.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	var res types.ContainerList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}
