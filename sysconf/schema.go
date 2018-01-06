/*
Copyright 2017 The Elasticshift Authors.
*/
package sysconf

import (
	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
)

func InitSchema(logger logrus.Logger, s Store) (queries graphql.Fields, mutations graphql.Fields) {

	r := &resolver{
		store:  s,
		logger: logger,
	}

	vcsSysconfFields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Represents the system config ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.VCSSysConf); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the system config",
		},

		"kind": &graphql.Field{
			Type:        graphql.String,
			Description: "Type of the system configuration",
		},

		"key": &graphql.Field{
			Type:        graphql.String,
			Description: "A key for the elasticshift application",
		},

		"secret": &graphql.Field{
			Type:        graphql.String,
			Description: "The secret for the elasticshift application",
		},

		"callback_url": &graphql.Field{
			Type:        graphql.String,
			Description: "The callback url for the elasticshift application",
		},
	}

	genericSysconfFields := graphql.Fields{

		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Represents the system config ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.GenericSysConf); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the system config",
		},

		"value": &graphql.Field{
			Type:        graphql.String,
			Description: "Value of the system configuration",
		},
	}

	accessModeEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "AccessMode",
		Values: graphql.EnumValueConfigMap{
			"ReadOnly": &graphql.EnumValueConfig{
				Value: 1,
			},

			"ReadWrite": &graphql.EnumValueConfig{
				Value: 2,
			},
		},
	})

	nfsSysconfFields := graphql.Fields{

		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Represents the system config ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.NFSVolumeSysConf); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"path": &graphql.Field{
			Type:        graphql.String,
			Description: "Path of the nfs share",
		},

		"server": &graphql.Field{
			Type:        graphql.String,
			Description: "NFS server address",
		},

		"accessmode": &graphql.Field{
			Type:        accessModeEnum,
			Description: "Access mode of the share",
		},
	}

	vcsSysconfType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "VCSSysconf",
			Fields:      vcsSysconfFields,
			Description: "An object of VCSSysConf type",
		},
	)

	genericSysconfType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "GenericSysconf",
			Fields:      genericSysconfFields,
			Description: "An object of GenericSysConf type",
		},
	)

	nfsSysconfType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "NFSSysconf",
			Fields:      nfsSysconfFields,
			Description: "An object of GenericSysConf type",
		},
	)

	queries = graphql.Fields{
		"VCSSysConf": &graphql.Field{
			Type: vcsSysconfType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Name of the system configuration",
				},
			},
			Resolve: r.FetchVCSSysConfByName,
		},

		"GenericSysConf": &graphql.Field{
			Type: genericSysconfType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the system configuration",
				},
			},
			Resolve: r.FetchGenericSysConfByName,
		},

		"NFSVolumeSysConf": &graphql.Field{
			Type: nfsSysconfType,
			Args: graphql.FieldConfigArgument{
				"server": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "NFS server host or ip",
				},
				"accessmode": &graphql.ArgumentConfig{
					Type:        accessModeEnum,
					Description: "Access mode such as ReadOnly or ReadWrite",
				},
			},
			Resolve: r.FetchNFSVolumeSysConfByName,
		},
	}

	mutations = graphql.Fields{
		"createVCSSysConf": &graphql.Field{
			Type: vcsSysconfType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the system configuration",
				},
				"key": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Key for the elasticshift oauth application",
				},
				"secret": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "secret for the elasticshift oauth application",
				},
				"callbackURL": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "callback url for the elasticshift oauth application",
				},
			},
			Resolve: r.CreateVCSSysConf,
		},

		"createGenericSysConf": &graphql.Field{
			Type: genericSysconfType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the system configuration",
				},

				"value": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Value for generic system configuration",
				},
			},
			Resolve: r.CreateGenericSysConf,
		},

		"createNFSVolumeSysConf": &graphql.Field{
			Type: nfsSysconfType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the nfs system configuration",
				},

				"server": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Server host or IP",
				},

				"accessmode": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(accessModeEnum),
					Description: "Access mode of the nfs share (ReadOnly | ReadWrite)",
				},
			},
			Resolve: r.CreateNFSVolumeSysConf,
		},
	}

	return queries, mutations
}
