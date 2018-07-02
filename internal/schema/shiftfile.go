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
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile"
	"gitlab.com/conspico/elasticshift/pkg/utils"
)

func newShiftfileSchema(ctx context.Context, logger logrus.Logger, s store.Shift) (queries graphql.Fields, mutations graphql.Fields) {

	r, _ := shiftfile.NewResolver(ctx, logger, s)

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Shiftfile identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.Shiftfile); ok {
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
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.Shiftfile); ok {
					return string(t.File), nil
				}
				return nil, nil
			},
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

	mutations = graphql.Fields{

		"addShiftfile": &graphql.Field{
			Type: shiftfileType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the shiftfile",
				},
				"description": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Description about the shiftfile",
				},
				"file": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The description file for elasticshift, which is a SHIFTFILE",
				},
				"team_id": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Team identifier",
				},
			},
			Resolve: r.AddShiftfile,
		},
	}

	return queries, mutations
}
