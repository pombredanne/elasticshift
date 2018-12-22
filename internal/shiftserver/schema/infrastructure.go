/*
Copyright 2018 The Elasticshift Authors.
*/
package schema

import (
	"context"

	"github.com/graphql-go/graphql"
	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	"github.com/elasticshift/elasticshift/internal/pkg/utils"
	"github.com/elasticshift/elasticshift/internal/shiftserver/infrastructure"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
)

func newInfrastructureSchema(ctx context.Context, loggr logger.Loggr, s store.Shift) (queries graphql.Fields, mutations graphql.Fields) {

	r, _ := infrastructure.NewResolver(ctx, loggr, s)

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Infrastructure identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.Infrastructure); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the infrastructure such as dev_gce_cluster",
		},

		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "Describe about the integration",
		},

		"kind": &graphql.Field{
			Type:        graphql.String,
			Description: "script type such as terraform, docker-compose or kubernetes",
		},

		"private": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "if true, only your team can see it, otherwise it's will be shown for other team as well (like sharing)",
		},

		"code": &graphql.Field{
			Type:        graphql.String,
			Description: "IAC - Infrastructure as code configuration",
		},

		"team": &graphql.Field{
			Type:        graphql.String,
			Description: "Team Identifier",
		},
	}

	infrastructureType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Infrastructure",
			Fields:      fields,
			Description: "An object of Container type",
		},
	)

	infrastructureArgs := graphql.FieldConfigArgument{

		"id": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Infrastructure identifier",
		},

		"team": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Team identifier",
		},
	}

	queries = graphql.Fields{
		"infrastructure": utils.MakeListType("InfrastructureList", infrastructureType, r.FetchInfrastructure, infrastructureArgs),
	}

	mutations = graphql.Fields{

		"addInfrastructure": &graphql.Field{
			Type: infrastructureType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the integration",
				},
				"description": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Describe about the infrastructe that it could create",
				},
				"kind": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Provider name such as onprem, gce, azure or amazon etc",
				},
				"private": &graphql.ArgumentConfig{
					Type:        graphql.Boolean,
					Description: "if true, only your team can see it, otherwise it''s will be shown for other team as well (like sharing)",
				},
				"code": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Script/config/code that brings up the infrastructure, Ex: docker-compose file, terraform script",
				},
				"team": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Team identifier",
				},
			},
			Resolve: r.AddInfrastructure,
		},
	}
	return queries, mutations
}
