/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/identity/team"
)

func InitSchema(logger logrus.FieldLogger, s Store, teamStore team.Store) (queries graphql.Fields, mutations graphql.Fields) {

	r := &resolver{
		store:     s,
		teamStore: teamStore,
		logger:    logger,
	}

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Represents the version control system ID",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.VCS); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the version control system, it shall be organization or user",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.VCS); ok {
					return t.Name, nil
				}
				return nil, nil
			},
		},

		"kind": &graphql.Field{
			Type:        graphql.String,
			Description: "Represents the repository type such as github, gitlab, bitbucket etc",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.VCS); ok {
					return t.Kind, nil
				}
				return nil, nil
			},
		},

		"ownerType": &graphql.Field{
			Type:        graphql.String,
			Description: "Represent the repository type sych as user or organization",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.VCS); ok {
					return t.OwnerType, nil
				}
				return nil, nil
			},
		},

		"avatarURL": &graphql.Field{
			Type:        graphql.String,
			Description: "An url that point the account profile picture",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.VCS); ok {
					return t.AvatarURL, nil
				}
				return nil, nil
			},
		},

		"accessToken": &graphql.Field{
			Type:        graphql.String,
			Description: "An access token that can be used to access this repository",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.VCS); ok {
					return t.AccessToken, nil
				}
				return nil, nil
			},
		},

		"refreshToken": &graphql.Field{
			Type:        graphql.String,
			Description: "The refresh token used to refresh the access token",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.VCS); ok {
					return t.RefreshToken, nil
				}
				return nil, nil
			},
		},

		"tokenExpiry": &graphql.Field{
			Type:        graphql.String,
			Description: "Time when the token will be expired",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.VCS); ok {
					return t.TokenExpiry.String(), nil
				}
				return nil, nil
			},
		},
	}

	vcsType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "VCS",
			Fields:      fields,
			Description: "An object of vcs type",
		},
	)

	queries = graphql.Fields{
		"vcs": &graphql.Field{
			Type: graphql.NewList(vcsType),
			Args: graphql.FieldConfigArgument{
				"teamID": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Represent the team ID",
				},
			},
			Resolve: r.FetchVCSByTeamID,
		},
	}

	mutations = graphql.Fields{}

	return queries, mutations
}
