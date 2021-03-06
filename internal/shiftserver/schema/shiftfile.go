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
	"github.com/elasticshift/elasticshift/internal/shiftserver/shiftfile"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
)

func newShiftfileSchema(ctx context.Context, loggr logger.Loggr, s store.Shift) (queries graphql.Fields, mutations graphql.Fields) {

	r, _ := shiftfile.NewResolver(ctx, loggr, s)

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
