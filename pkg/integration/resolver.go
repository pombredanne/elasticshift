/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errIDCantBeEmpty     = errors.New("Integration ID cannot be empty")
	errTeamCannotBeEmpty = errors.New("Team must be provided")
)

type resolver struct {
	store  Store
	logger logrus.Logger
	Ctx    context.Context
}

func (r *resolver) FetchIntegration(params graphql.ResolveParams) (interface{}, error) {

	team := params.Args["team"].(string)
	if team == "" {
		return nil, errTeamCannotBeEmpty
	}

	q := bson.M{"team": team}

	id := params.Args["id"].(string)
	if id != "" {
		q["_id"] = bson.ObjectIdHex(id)
	}

	var err error
	var result []types.Integration
	r.store.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	var res types.IntegrationList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}

func (r *resolver) AddKubernetesCluster(params graphql.ResolveParams) (interface{}, error) {

	team, _ := params.Args["team"].(string)
	kind, _ := params.Args["kind"].(int)
	provider, _ := params.Args["provider"].(int)
	host, _ := params.Args["host"].(string)
	certificate, _ := params.Args["certificate"].(string)
	token, _ := params.Args["token"].(string)
	name, _ := params.Args["name"].(string)

	i := &types.Integration{}
	i.Name = name
	i.Team = team
	i.Kind = kind
	i.Host = host
	i.Certificate = certificate
	i.Token = token
	i.Provider = provider

	err := r.store.Save(i)
	if err != nil {
		return nil, fmt.Errorf("Failed to add integration: %v", err)
	}
	return i, nil
}
