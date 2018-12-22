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
	"github.com/elasticshift/elasticshift/internal/shiftserver/repository"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
)

func newRepositorySchema(
	ctx context.Context,
	loggr logger.Loggr,
	s store.Shift,
) (queries graphql.Fields, mutations graphql.Fields) {

	r, _ := repository.NewResolver(ctx, loggr, s)

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Repository identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.Repository); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"vcs_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Version control system identifier",
		},

		"repo_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier belong to the repository in target scm",
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the repository or project",
		},

		"private": &graphql.Field{
			Type:        graphql.String,
			Description: "True if it is a private repository, otherwise public",
		},

		"link": &graphql.Field{
			Type:        graphql.String,
			Description: "Hyperlink to the repository",
		},

		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "Additional information about the repository",
		},

		"fork": &graphql.Field{
			Type:        graphql.String,
			Description: "True if this is a forked repo",
		},

		"default_branch": &graphql.Field{
			Type:        graphql.String,
			Description: "Defualt branch for this repository",
		},

		"language": &graphql.Field{
			Type:        graphql.String,
			Description: "Represent the source code language reside in this repository",
		},

		"source": &graphql.Field{
			Type:        graphql.String,
			Description: "The source of the repository such as github.com, gitlab.com etc",
		},

		"build": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "builds",
				Fields: graphql.Fields{
					"nodes": &graphql.Field{Type: graphql.NewList(BuildType)},
					"count": &graphql.Field{Type: graphql.Int},
				},
			}),
			Resolve: r.FetchBuild,
		},
	}

	repositoryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Repository",
			Fields:      fields,
			Description: "An object of repository type",
		},
	)

	args := graphql.FieldConfigArgument{
		"team": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Represent the team name",
		},

		"vcs_id": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "Represent the version control system identifier",
		},
	}

	queries = graphql.Fields{
		"repository": utils.MakeListType("RepositoryList", repositoryType, r.FetchRepository, args),
	}

	mutations = graphql.Fields{

		"addRepository": &graphql.Field{
			Type: repositoryType,
			Args: graphql.FieldConfigArgument{
				"uri": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "URI of the repository that typically used to clone. Ex: git@github.com:<account>/<project>.git",
				},
				"team": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Represent the team name or ID",
				},
			},
			Resolve: r.AddRepository,
		},
	}

	return queries, mutations
}
