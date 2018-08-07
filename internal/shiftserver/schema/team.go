/*
Copyright 2017 The Elasticshift Authors.
*/
package schema

import (
	"context"

	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/team"
)

func newTeamSchema(ctx context.Context, loggr logger.Loggr, s store.Shift) (queries graphql.Fields, mutations graphql.Fields) {

	r, _ := team.NewResolver(ctx, loggr, s)

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Represents the team ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the team",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.Name, nil
				}
				return nil, nil
			},
		},

		"display": &graphql.Field{
			Type:        graphql.String,
			Description: "Name that is used to represent for display purpose such as logged in name etc",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.Display, nil
				}
				return nil, nil
			},
		},

		"accounts": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "List of version control system accounts linked for this team",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Team); ok {
					return t.Accounts, nil
				}
				return nil, nil
			},
		},
	}

	teamType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Team",
			Fields:      fields,
			Description: "An object of team type",
		},
	)

	queries = graphql.Fields{
		"team": &graphql.Field{
			Type: teamType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Represent the team ID",
				},
				"name": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Name of the team",
				},
			},
			Resolve: r.FetchByNameOrID,
		},
	}

	mutations = graphql.Fields{
		"createTeam": &graphql.Field{
			Type: teamType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the team",
				},
			},
			Resolve: r.CreateTeam,
		},
	}

	return queries, mutations
}
