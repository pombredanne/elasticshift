/*
Copyright 2018 The Elasticshift Authors.
*/
package app

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/pkg/utils"
)

func InitSchema(logger logrus.Logger, ctx context.Context, s Store) (queries graphql.Fields, mutations graphql.Fields) {

	r := &resolver{
		store:  s,
		logger: logger,
		Ctx:    ctx,
	}

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "App identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.App); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the app",
		},

		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "Description of the app",
		},

		"language": &graphql.Field{
			Type:        graphql.String,
			Description: "Language used to develop this app",
		},

		"version": &graphql.Field{
			Type:        graphql.String,
			Description: "Version of the app",
		},

		"used_team_count": &graphql.Field{
			Type:        graphql.Int,
			Description: "Number of team currently using this app",
		},

		"used_build_count": &graphql.Field{
			Type:        graphql.Int,
			Description: "Total build using this app",
		},

		"icon_url": &graphql.Field{
			Type:        graphql.String,
			Description: "Time when the container destroyed",
		},

		"source_url": &graphql.Field{
			Type:        graphql.String,
			Description: "Time when the container destroyed",
		},

		"readme": &graphql.Field{
			Type:        graphql.String,
			Description: "Time when the container destroyed",
		},

		"ratings": &graphql.Field{
			Type:        graphql.String,
			Description: "Time when the container destroyed",
		},
	}

	appType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "App",
			Fields:      fields,
			Description: "An object of App type",
		},
	)

	appArgs := graphql.FieldConfigArgument{
		"team": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Team identifier",
		},

		"name": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "App name",
		},
	}

	queries = graphql.Fields{
		"app": utils.MakeListType("AppList", appType, r.FetchApp, appArgs),
	}

	mutations = graphql.Fields{}

	return queries, mutations
}
