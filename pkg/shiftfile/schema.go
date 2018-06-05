/*
Copyright 2018 The Elasticshift Authors.
*/
package shiftfile

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
			Description: "Shiftfile identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.Plugin); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the shiftfile",
		},

		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "Description of the shiftfile",
		},

		"file": &graphql.Field{
			Type:        graphql.String,
			Description: "The shiftfile",
		},

		"used_by_repos": &graphql.Field{
			Type:        graphql.Int,
			Description: "Number of repo(s) currently using this shiftfile",
		},

		"used_by_teams": &graphql.Field{
			Type:        graphql.Int,
			Description: "Number of team(s) currently using this shiftfile",
		},

		"ratings": &graphql.Field{
			Type:        graphql.String,
			Description: "The ratings of the shiftfile",
		},
	}

	shiftfileType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Shiftfile",
			Fields:      fields,
			Description: "An object of Shiftfile type",
		},
	)

	shiftfileArgs := graphql.FieldConfigArgument{

		"name": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "Name of the shiftfile",
		},
	}

	queries = graphql.Fields{
		"shiftfile": utils.MakeListType("ShiftfileList", shiftfileType, r.FetchShiftfile, shiftfileArgs),
	}

	mutations = graphql.Fields{}

	return queries, mutations
}
