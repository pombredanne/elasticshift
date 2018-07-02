/*
Copyright 2018 The Elasticshift Authors.
*/
package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errTeamCannotBeEmpty = errors.New("Team must be provided")
)

// Resolver ...
type Resolver interface {
	FetchInfrastructure(params graphql.ResolveParams) (interface{}, error)
	AddInfrastructure(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	store  store.Infrastructure
	logger logrus.Logger
	Ctx    context.Context
}

// NewResolver ...
func NewResolver(ctx context.Context, logger logrus.Logger, s store.Shift) (Resolver, error) {

	r := &resolver{
		store:  s.Infrastructure,
		logger: logger,
		Ctx:    ctx,
	}
	return r, nil
}

func (r *resolver) FetchInfrastructure(params graphql.ResolveParams) (interface{}, error) {

	team := params.Args["team"].(string)
	if team == "" {
		return nil, errTeamCannotBeEmpty
	}

	q := bson.M{"team": team}

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

func (r *resolver) AddInfrastructure(params graphql.ResolveParams) (interface{}, error) {

	team, _ := params.Args["team"].(string)
	name, _ := params.Args["name"].(string)
	description, _ := params.Args["description"].(string)
	kind, _ := params.Args["kind"].(string)
	private, _ := params.Args["private"].(bool)
	code, _ := params.Args["code"].(string)

	i := &types.Infrastructure{}
	i.Name = name
	i.Description = description
	i.Kind = kind
	i.Private = private
	i.Code = code
	i.Team = team

	err := r.store.Save(i)
	if err != nil {
		return nil, fmt.Errorf("Failed to add infrastructure : %v", err)
	}

	return i, err
}
