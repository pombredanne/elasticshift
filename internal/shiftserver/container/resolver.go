/*
Copyright 2018 The Elasticshift Authors.
*/
package container

import (
	"context"
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errIDCantBeEmpty     = errors.New("Container ID cannot be empty")
	errTeamCannotBeEmpty = errors.New("Team must be provided")
)

// Resolver ...
type Resolver interface {
	FetchContainer(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	store  store.Container
	logger *logrus.Entry
	Ctx    context.Context
}

// NewResolver ...
func NewResolver(ctx context.Context, loggr logger.Loggr, s store.Shift) (Resolver, error) {

	r := &resolver{
		store:  s.Container,
		logger: loggr.GetLogger("graphql/container"),
		Ctx:    ctx,
	}
	return r, nil
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
