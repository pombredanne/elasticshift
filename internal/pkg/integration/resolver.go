/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errIDCantBeEmpty     = errors.New("Integration ID cannot be empty")
	errTeamCannotBeEmpty = errors.New("Team must be provided")
)

const (
	INT_ContainerEngine int = iota + 1
	INT_Storage
)

// Resolver ...
type Resolver interface {
	FetchContainerEngine(params graphql.ResolveParams) (interface{}, error)
	FetchStorage(params graphql.ResolveParams) (interface{}, error)
	AddKubernetesCluster(params graphql.ResolveParams) (interface{}, error)
	AddStorage(params graphql.ResolveParams) (interface{}, error) 
}

type resolver struct {
	store  store.Integration
	logger logrus.Logger
	Ctx    context.Context
}

// NewResolver ...
func NewResolver(ctx context.Context, logger logrus.Logger, s store.Shift) (Resolver, error) {

	r := &resolver{
		store:  s.Integration,
		logger: logger,
		Ctx:    ctx,
	}
	return r, nil
}

func (r *resolver) FetchContainerEngine(params graphql.ResolveParams) (interface{}, error) {

	team, _ := params.Args["team"].(string)
	if team == "" {
		return nil, errTeamCannotBeEmpty
	}

	q := bson.M{"team": team, "internal_type": INT_ContainerEngine}

	id, _ := params.Args["id"].(string)
	if id != "" {
		q["_id"] = bson.ObjectIdHex(id)
	}

	var err error
	var result []types.ContainerEngine
	r.store.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	var res types.ContainerEngineList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}

func (r *resolver) FetchStorage(params graphql.ResolveParams) (interface{}, error) {

	team, _ := params.Args["team"].(string)
	if team == "" {
		return nil, errTeamCannotBeEmpty
	}

	q := bson.M{"team": team, "internal_type": INT_Storage}

	id, _ := params.Args["id"].(string)
	if id != "" {
		q["_id"] = bson.ObjectIdHex(id)
	}

	var err error
	var result []types.Storage
	r.store.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	var res types.StorageList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}

func (r *resolver) AddKubernetesCluster(params graphql.ResolveParams) (interface{}, error) {

	team, _ := params.Args["team"].(string)
	name, _ := params.Args["name"].(string)

	var ce types.ContainerEngine
	err := r.store.FindOne(bson.M{"team": team, "name": name}, &ce)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return nil, fmt.Errorf("Failed to check if the given kubernetes integration already exist :%v", err)
	}

	if ce.ID.Hex() != "" {
		return nil, fmt.Errorf("The container engine name '%s' already exist for your team", name)
	}

	kind, _ := params.Args["kind"].(int)
	provider, _ := params.Args["provider"].(int)
	host, _ := params.Args["host"].(string)
	certificate, _ := params.Args["certificate"].(string)
	token, _ := params.Args["token"].(string)

	i := types.ContainerEngine{}
	i.Name = name
	i.Team = team
	i.Kind = kind
	i.Host = host
	i.Certificate = certificate
	i.Token = token
	i.Provider = provider
	i.InternalType = INT_ContainerEngine

	err = r.store.Save(&i)
	if err != nil {
		return nil, fmt.Errorf("Failed to add integration: %v", err)
	}
	return i, nil
}

func (r *resolver) AddStorage(params graphql.ResolveParams) (interface{}, error) {

	name, _ := params.Args["name"].(string)
	team, _ := params.Args["team"].(string)

	var stor types.Storage
	err := r.store.FindOne(bson.M{"team": team, "name": name}, &stor)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return nil, fmt.Errorf("Failed to check if the given storage integration already exist :%v", err)
	}

	if stor.ID.Hex() != "" {
		return nil, fmt.Errorf("The storage name '%s' already exist for your team", name)
	}

	kind, _ := params.Args["kind"].(int)
	provider, _ := params.Args["provider"].(int)
	host, _ := params.Args["host"].(string)
	certificate, _ := params.Args["certificate"].(string)
	accesskey, _ := params.Args["accesskey"].(string)
	secretkey, _ := params.Args["secretkey"].(string)

	i := types.Storage{}
	i.Name = name
	i.Team = team
	i.Kind = kind
	i.Host = host
	i.Certificate = certificate
	i.AccessKey = accesskey
	i.SecretKey = secretkey
	i.Provider = provider
	i.InternalType = INT_Storage

	err = r.store.Save(&i)
	if err != nil {
		return nil, fmt.Errorf("Failed to add integration: %v", err)
	}
	return i, nil
}
