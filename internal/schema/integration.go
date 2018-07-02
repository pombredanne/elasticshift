/*
Copyright 2018 The Elasticshift Authors.
*/
package schema

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/store"
	"gitlab.com/conspico/elasticshift/internal/pkg/integration"
	"gitlab.com/conspico/elasticshift/pkg/utils"
)

func newIntegrationSchema(ctx context.Context, logger logrus.Logger, s store.Shift) (queries graphql.Fields, mutations graphql.Fields) {

	r, _ := integration.NewResolver(ctx, logger, s)

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
		Name: "ContainerEngineKind",
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

	storageKindEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "StorageKind",
		Values: graphql.EnumValueConfigMap{

			"MINIO": &graphql.EnumValueConfig{
				Value: 1,
			},

			"S3": &graphql.EnumValueConfig{
				Value: 2,
			},

			"GCE": &graphql.EnumValueConfig{
				Value: 3,
			},
		},
	})

	containerEngineFields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Integration identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.ContainerEngine); ok {
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

				if t, ok := p.Source.(types.ContainerEngine); ok {
					return int(t.Kind), nil
				}
				return nil, nil
			},
		},

		"provider": &graphql.Field{
			Type:        providerEnum,
			Description: "Provider such as ONPREM, GCE, AMAZON, ALIBABACLOUD and azure etc",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				if t, ok := p.Source.(types.ContainerEngine); ok {
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

	storageFields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Integration identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.Storage); ok {
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
			Type:        storageKindEnum,
			Description: "Type of storage such as minio, gce or amazon s3",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				if t, ok := p.Source.(types.Storage); ok {
					return int(t.Kind), nil
				}
				return nil, nil
			},
		},

		"provider": &graphql.Field{
			Type:        providerEnum,
			Description: "Provider such as ONPREM, GCE, AMAZON, ALIBABACLOUD and azure etc",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				if t, ok := p.Source.(types.Storage); ok {
					return int(t.Provider), nil
				}
				return nil, nil
			},
		},

		"certificate": &graphql.Field{
			Type:        graphql.String,
			Description: "Certificate to access the storage provider",
		},

		"host": &graphql.Field{
			Type:        graphql.String,
			Description: "Host of the storage provider",
		},

		"access_key": &graphql.Field{
			Type:        graphql.String,
			Description: "Token to access the storage provider",
		},

		"secret_key": &graphql.Field{
			Type:        graphql.String,
			Description: "Token to access the storage provider",
		},

		"team": &graphql.Field{
			Type:        graphql.String,
			Description: "Team Identifier",
		},
	}

	containerEngineType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "ContainerEngine",
			Fields:      containerEngineFields,
			Description: "An object of ContainerEngine type",
		},
	)

	storageType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Storage",
			Fields:      storageFields,
			Description: "An object of Storage type",
		},
	)

	containerEngineArgs := graphql.FieldConfigArgument{

		"id": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "Container Engine identifier",
		},

		"team": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Team identifier",
		},
	}

	storageArgs := graphql.FieldConfigArgument{

		"id": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "Storage identifier",
		},

		"team": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Team identifier",
		},
	}

	queries = graphql.Fields{
		"ContainerEngine": utils.MakeListType("ContainerEngineList", containerEngineType, r.FetchContainerEngine, containerEngineArgs),
		"Storage":         utils.MakeListType("StorageList", storageType, r.FetchStorage, storageArgs),
	}

	mutations = graphql.Fields{

		"addKubernetesCluster": &graphql.Field{
			Type: containerEngineType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the container engine",
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

		"addStorage": &graphql.Field{
			Type: storageType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the storage",
				},
				"kind": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(storageKindEnum),
					Description: "Type of provider where the storage cluster resides such as kubernetes, dcos or dockerswarm",
				},
				"provider": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(providerEnum),
					Description: "Provider name such as onprem, gce, azure or amazon etc",
				},
				"host": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Host of the storage provider",
				},
				"certificate": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Certificate for the storage provider",
				},
				"accesskey": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Access key to storage provider",
				},
				"secretkey": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "secret key to storage provider",
				},
				"team": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Team identifier",
				},
			},
			Resolve: r.AddStorage,
		},
	}

	return queries, mutations
}
