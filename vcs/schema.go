/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/core/utils"
	"gitlab.com/conspico/elasticshift/identity/oauth2/providers"
	"gitlab.com/conspico/elasticshift/identity/team"
	"gitlab.com/conspico/elasticshift/vcs/repository"
)

func InitSchema(logger logrus.Logger, providers providers.Providers, s Store, teamStore team.Store, repositoryStore repository.Store) (queries graphql.Fields, mutations graphql.Fields) {

	r := &resolver{
		store:           s,
		repositoryStore: repositoryStore,
		teamStore:       teamStore,
		logger:          logger,
		providers:       providers,
	}

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Represents the version control system ID",
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the version control system, it shall be organization or user",
		},

		"kind": &graphql.Field{
			Type:        graphql.String,
			Description: "Represents the repository type such as github, gitlab, bitbucket etc",
		},

		"owner_type": &graphql.Field{
			Type:        graphql.String,
			Description: "Represent the repository type sych as user or organization",
		},

		"avatar": &graphql.Field{
			Type:        graphql.String,
			Description: "An url that point the account profile picture",
		},

		"access_token": &graphql.Field{
			Type:        graphql.String,
			Description: "An access token that can be used to access this repository",
		},

		"refresh_token": &graphql.Field{
			Type:        graphql.String,
			Description: "The refresh token used to refresh the access token",
		},

		"token_expiry": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "Time when the token will be expired",
		},
	}

	repositoryFields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Repository identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Repository); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"vcs_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Version control system identifier",
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
	}

	vcsType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "VCS",
			Fields:      fields,
			Description: "An object of vcs type",
		},
	)

	repositoryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Repository",
			Fields:      repositoryFields,
			Description: "An object of repository type",
		},
	)

	teamArg := graphql.FieldConfigArgument{
		"team": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Represent the team name or ID",
		},
	}

	queries = graphql.Fields{
		"repository": utils.MakeListType("RepositoryList", repositoryType, r.FetchRepository, teamArg),
		"vcs":        utils.MakeListType("VCSList", vcsType, r.FetchVCS, teamArg),
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
