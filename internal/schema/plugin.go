/*
Copyright 2018 The Elasticshift Authors.
*/
package schema

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/pkg/plugin"
	"gitlab.com/conspico/elasticshift/internal/store"
	"gitlab.com/conspico/elasticshift/pkg/utils"
)

func newPluginSchema(ctx context.Context, logger logrus.Logger, s store.Shift) (queries graphql.Fields, mutations graphql.Fields) {

	r, _ := plugin.NewResolver(ctx, logger, s)

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Plugin identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.Plugin); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the plugin",
		},

		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "Description of the plugin",
		},

		"language": &graphql.Field{
			Type:        graphql.String,
			Description: "Language used to develop the plugin",
		},

		"version": &graphql.Field{
			Type:        graphql.String,
			Description: "Version of the plugin",
		},

		"used_team_count": &graphql.Field{
			Type:        graphql.Int,
			Description: "Number of team(s) currently using the plugin",
		},

		"used_build_count": &graphql.Field{
			Type:        graphql.Int,
			Description: "Total number of build(s) using the plugin",
		},

		"icon_url": &graphql.Field{
			Type:        graphql.String,
			Description: "Plugin icon url",
		},

		"source_url": &graphql.Field{
			Type:        graphql.String,
			Description: "Source code reference of the plugin",
		},

		"readme": &graphql.Field{
			Type:        graphql.String,
			Description: "Document that refers the readme",
		},

		"ratings": &graphql.Field{
			Type:        graphql.String,
			Description: "Rating of the plugin",
		},
	}

	pluginType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Plugin",
			Fields:      fields,
			Description: "An object of Plugin type",
		},
	)

	pluginArgs := graphql.FieldConfigArgument{
		"team": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Team identifier",
		},

		"name": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "Plugin name",
		},
	}

	queries = graphql.Fields{
		"plugin": utils.MakeListType("PluginList", pluginType, r.FetchPlugin, pluginArgs),
	}

	mutations = graphql.Fields{}

	return queries, mutations
}
