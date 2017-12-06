/*
Copyright 2017 The Elasticshift Authors.
*/
package team

import (
	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	core "gitlab.com/conspico/elasticshift/core/store"
	"gitlab.com/conspico/elasticshift/core/utils"
)

func InitSchema(s core.Store, logger logrus.FieldLogger) (queries graphql.Fields, mutations graphql.Fields) {

	r := &resolver{
		store:  NewStore(s),
		logger: logger,
	}

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.Name, nil
				}
				return nil, nil
			},
		},

		"display": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.Display, nil
				}
				return nil, nil
			},
		},

		"accounts": &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.Accounts, nil
				}
				return nil, nil
			},
		},
	}

	teamType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Team",
			Fields: fields,
		},
	)

	queries = graphql.Fields{
		"team": &graphql.Field{
			Type: teamType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: r.FetchByNameOrID,
		},
		"teams": utils.MakeListField(utils.MakeNodeListType("TeamList", teamType), r.FetchTeams),
	}

	mutations = graphql.Fields{
		"createTeam": &graphql.Field{
			Type: teamType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: r.CreateTeam,
		},
	}

	return queries, mutations
}
