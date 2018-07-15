/*
Copyright 2018 The Elasticshift Authors.
*/
package schema

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/integration"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
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

			"AmazonS3": &graphql.EnumValueConfig{
				Value: 2,
			},

			"GoogleCloudStorage": &graphql.EnumValueConfig{
				Value: 3,
			},

			"NFS": &graphql.EnumValueConfig{
				Value: 4,
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

	minioFields := graphql.Fields{

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
	}

	minioStorageType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "MinioStorage",
			Fields:      minioFields,
			Description: "An object of Minio Storage type",
		},
	)

	nfsFields := graphql.Fields{

		"server": &graphql.Field{
			Type:        graphql.String,
			Description: "Host of the nfs server",
		},

		"path": &graphql.Field{
			Type:        graphql.String,
			Description: "Mount of the nfs server",
		},

		"readonly": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Indicates the mount is readonly",
		},

		"mount_path": &graphql.Field{
			Type:        graphql.String,
			Description: "Indicates the mount path while mounting to container as storage",
		},
	}

	nfsStorageType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "NFSStorage",
			Fields:      nfsFields,
			Description: "An object of NFS Storage type",
		},
	)

	sourceFields := graphql.Fields{
		"nfs": &graphql.Field{
			Type:        nfsStorageType,
			Description: "NFS Storage",
		},

		"minio": &graphql.Field{
			Type:        minioStorageType,
			Description: "Minio Storage",
		},
	}

	sourceType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "StorageSource",
			Fields:      sourceFields,
			Description: "An object of StorageSource type",
		},
	)

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

		"team": &graphql.Field{
			Type:        graphql.String,
			Description: "Team Identifier",
		},

		"storage_source": &graphql.Field{
			Type:        sourceType,
			Description: "Source of the storage such as nfs, minio etc",
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

	nfsInputType := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "NFSStorageInput",
			Fields: graphql.InputObjectConfigFieldMap{

				"server": &graphql.InputObjectFieldConfig{
					Type:        graphql.String,
					Description: "Host of the nfs server",
				},

				"path": &graphql.InputObjectFieldConfig{
					Type:        graphql.String,
					Description: "Mount of the nfs server",
				},

				"readonly": &graphql.InputObjectFieldConfig{
					Type:        graphql.Boolean,
					Description: "Indicates the mount is readonly",
				},

				"mount_path": &graphql.InputObjectFieldConfig{
					Type:        graphql.String,
					Description: "Mount of the nfs server when mounting to container (shift_dir)",
				},
			},
		},
	)

	minioInputType := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "MinioStorageInput",
			Fields: graphql.InputObjectConfigFieldMap{

				"host": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Host of the storage provider",
				},
				"certificate": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Certificate for the storage provider",
				},
				"accesskey": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Access key to storage provider",
				},
				"secretkey": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "secret key to storage provider",
				},
			},
		},
	)

	storageSourceInputType := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "StorageSourceInput",
			Fields: graphql.InputObjectConfigFieldMap{
				"minio": &graphql.InputObjectFieldConfig{
					Type:        minioInputType,
					Description: "Minio Input",
				},

				"nfs": &graphql.InputObjectFieldConfig{
					Type:        nfsInputType,
					Description: "NFS Input",
				},
			},
		},
	)

	storageInputType := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "StorageInput",
			Fields: graphql.InputObjectConfigFieldMap{

				"name": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the storage",
				},
				"kind": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(storageKindEnum),
					Description: "Type of provider where the storage cluster resides such as kubernetes, dcos or dockerswarm",
				},
				"provider": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(providerEnum),
					Description: "Provider name such as onprem, gce, azure or amazon etc",
				},
				"team": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Team identifier",
				},
				"storage_source": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(storageSourceInputType),
					Description: "Storage source Input",
				},
			},
		},
	)

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
				"storage": &graphql.ArgumentConfig{
					Type:        storageInputType,
					Description: "Storage Input",
				},
				// "name": &graphql.ArgumentConfig{
				// 	Type:        graphql.NewNonNull(graphql.String),
				// 	Description: "Name of the storage",
				// },
				// "kind": &graphql.ArgumentConfig{
				// 	Type:        graphql.NewNonNull(storageKindEnum),
				// 	Description: "Type of provider where the storage cluster resides such as kubernetes, dcos or dockerswarm",
				// },
				// "provider": &graphql.ArgumentConfig{
				// 	Type:        graphql.NewNonNull(providerEnum),
				// 	Description: "Provider name such as onprem, gce, azure or amazon etc",
				// },
				// "host": &graphql.ArgumentConfig{
				// 	Type:        graphql.NewNonNull(graphql.String),
				// 	Description: "Host of the storage provider",
				// },
				// "certificate": &graphql.ArgumentConfig{
				// 	Type:        graphql.NewNonNull(graphql.String),
				// 	Description: "Certificate for the storage provider",
				// },
				// "accesskey": &graphql.ArgumentConfig{
				// 	Type:        graphql.NewNonNull(graphql.String),
				// 	Description: "Access key to storage provider",
				// },
				// "secretkey": &graphql.ArgumentConfig{
				// 	Type:        graphql.NewNonNull(graphql.String),
				// 	Description: "secret key to storage provider",
				// },
				// "team": &graphql.ArgumentConfig{
				// 	Type:        graphql.NewNonNull(graphql.String),
				// 	Description: "Team identifier",
				// },
			},
			Resolve: r.AddStorage,
		},
	}

	return queries, mutations
}
