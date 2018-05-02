/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

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

	providerEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "Provider",
		Values: graphql.EnumValueConfigMap{
			"ONPREM": &graphql.EnumValueConfig{
				Value: 1,
			},

			"GCE": &graphql.EnumValueConfig{
				Value: 2,
			},

			"AZURE": &graphql.EnumValueConfig{
				Value: 3,
			},

			"AMAZON": &graphql.EnumValueConfig{
				Value: 4,
			},

			"DIGITALOCEAN": &graphql.EnumValueConfig{
				Value: 5,
			},

			"ALIBABACLOUD": &graphql.EnumValueConfig{
				Value: 6,
			},
		},
	})

	kindEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "Kind",
		Values: graphql.EnumValueConfigMap{
			"KUBERNETES": &graphql.EnumValueConfig{
				Value: 1,
			},

			"DOCKERSWARM": &graphql.EnumValueConfig{
				Value: 2,
			},

			"DCOS": &graphql.EnumValueConfig{
				Value: 3,
			},
		},
	})

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Integration identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.Integration); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the integration such as dev_gce_cluster",
		},

		"kind": &graphql.Field{
			Type:        kindEnum,
			Description: "Type of cluster such as kubernetes, dcos or dockerswarm",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				if t, ok := p.Source.(types.Integration); ok {
					return int(t.Kind), nil
				}
				return nil, nil
			},
		},

		"provider": &graphql.Field{
			Type:        providerEnum,
			Description: "Provider such as ONPREM, GCE, AMAZON, ALIBABACLOUD and azure etc",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				if t, ok := p.Source.(types.Integration); ok {
					return int(t.Provider), nil
				}
				return nil, nil
			},
		},

		"certificate": &graphql.Field{
			Type:        graphql.String,
			Description: "Certificate to access the cluster",
		},

		"host": &graphql.Field{
			Type:        graphql.String,
			Description: "Host of the cluster",
		},

		"token": &graphql.Field{
			Type:        graphql.String,
			Description: "Token to access the cluster",
		},

		"team": &graphql.Field{
			Type:        graphql.String,
			Description: "Team Identifier",
		},
	}

	integrationType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Integration",
			Fields:      fields,
			Description: "An object of Integration type",
		},
	)

	integrationArgs := graphql.FieldConfigArgument{

		"id": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "Integration identifier",
		},

		"team": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Team identifier",
		},
	}

	queries = graphql.Fields{
		"integration": utils.MakeListType("IntegrationList", integrationType, r.FetchIntegration, integrationArgs),
	}

	mutations = graphql.Fields{

		"addIntegration": &graphql.Field{
			Type: integrationType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the integration",
				},
				"kind": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(kindEnum),
					Description: "Type of cluster such as kubernetes, dcos or dockerswarm",
				},
				"provider": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(providerEnum),
					Description: "Provider name such as onprem, gce, azure or amazon etc",
				},
				"host": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Host of the cluster",
				},
				"certificate": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Certificate for the cluster",
				},
				"token": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Token for the cluster",
				},
				"team": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Team identifier",
				},
			},
			Resolve: r.AddKubernetesCluster,
		},
	}

	return queries, mutations
}
