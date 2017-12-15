/*
Copyright 2017 The Elasticshift Authors.
*/
package sysconf

import (
	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
)

func InitSchema(logger logrus.FieldLogger, s Store) (queries graphql.Fields, mutations graphql.Fields) {

	r := &resolver{
		store:  s,
		logger: logger,
	}

	vcsSysconfFields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Represents the system config ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the system config",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.Name, nil
				}
				return nil, nil
			},
		},

		"kind": &graphql.Field{
			Type:        graphql.String,
			Description: "Type of the system configuration",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.Display, nil
				}
				return nil, nil
			},
		},

		"key": &graphql.Field{
			Type:        graphql.String,
			Description: "A key for the elasticshift application",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.Accounts, nil
				}
				return nil, nil
			},
		},

		"secret": &graphql.Field{
			Type:        graphql.String,
			Description: "The secret for the elasticshift application",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.Accounts, nil
				}
				return nil, nil
			},
		},

		"callbackURL": &graphql.Field{
			Type:        graphql.String,
			Description: "The callback url for the elasticshift application",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.Accounts, nil
				}
				return nil, nil
			},
		},
	}

	vcsSysconfType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "VCSSysconf",
			Fields:      vcsSysconfFields,
			Description: "An object of VCSSysConf type",
		},
	)

	queries = graphql.Fields{
		"vcs": &graphql.Field{
			Type: vcsSysconfType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Name of the system configuration",
				},
			},
			Resolve: r.FetchVCSSysConfByName,
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
	}

	return queries, mutations
}
