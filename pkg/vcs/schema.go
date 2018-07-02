/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/internal/store"
	"gitlab.com/conspico/elasticshift/pkg/identity/oauth2/providers"
	"gitlab.com/conspico/elasticshift/pkg/utils"
)

func InitSchema(
	logger logrus.Logger,
	providers providers.Providers,
	s store.Shift,
) (queries graphql.Fields, mutations graphql.Fields) {

	r := &resolver{
		store:     s.Vcs,
		teamStore: s.Team,
		logger:    logger,
		providers: providers,
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

	vcsType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "VCS",
			Fields:      fields,
			Description: "An object of vcs type",
		},
	)

	args := graphql.FieldConfigArgument{
		"team": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Represent the team name or ID",
		},
	}

	queries = graphql.Fields{
		"vcs": utils.MakeListType("VCSList", vcsType, r.FetchVCS, args),
	}

	mutations = graphql.Fields{}

	return queries, mutations
}
