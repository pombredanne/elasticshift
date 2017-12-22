/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/identity/team"
)

func InitSchema(logger logrus.Logger, s Store, teamStore team.Store) (queries graphql.Fields, mutations graphql.Fields) {

	r := &resolver{
		store:     s,
		teamStore: teamStore,
		logger:    logger,
	}

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Represents the version control system ID",
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 	if t, ok := p.Source.(*types.VCS); ok {
			// 		logger.Warnln("Type from VCS ID field: ", t)
			// 		return t.ID, nil
			// 	}
			// 	return nil, nil
			// },
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the version control system, it shall be organization or user",
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 	if t, ok := p.Source.(*types.VCS); ok {
			// 		return t.Name, nil
			// 	}
			// 	return nil, nil
			// },
		},

		"kind": &graphql.Field{
			Type:        graphql.String,
			Description: "Represents the repository type such as github, gitlab, bitbucket etc",
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 	if t, ok := p.Source.(*types.VCS); ok {
			// 		return strconv.Itoa(t.Kind), nil
			// 	}
			// 	return nil, nil
			// },
		},

		"owner_type": &graphql.Field{
			Type:        graphql.String,
			Description: "Represent the repository type sych as user or organization",
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 	if t, ok := p.Source.(*types.VCS); ok {
			// 		return t.OwnerType, nil
			// 	}
			// 	return nil, nil
			// },
		},

		"avatar": &graphql.Field{
			Type:        graphql.String,
			Description: "An url that point the account profile picture",
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 	if t, ok := p.Source.(*types.VCS); ok {
			// 		return t.AvatarURL, nil
			// 	}
			// 	return nil, nil
			// },
		},

		"access_token": &graphql.Field{
			Type:        graphql.String,
			Description: "An access token that can be used to access this repository",
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 	if t, ok := p.Source.(*types.VCS); ok {
			// 		return t.AccessToken, nil
			// 	}
			// 	return nil, nil
			// },
		},

		"refresh_token": &graphql.Field{
			Type:        graphql.String,
			Description: "The refresh token used to refresh the access token",
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 	if t, ok := p.Source.(*types.VCS); ok {
			// 		return t.RefreshToken, nil
			// 	}
			// 	return nil, nil
			// },
		},

		"token_expiry": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "Time when the token will be expired",
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 	if t, ok := p.Source.(*types.VCS); ok {
			// 		return t.TokenExpiry.String(), nil
			// 	}
			// 	return nil, nil
			// },
		},
	}

	vcsType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "VCS",
			Fields:      fields,
			Description: "An object of vcs type",
		},
	)

	listType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "VCSList",
			Fields: graphql.Fields{
				"nodes": &graphql.Field{
					Type: graphql.NewList(vcsType),
				},
				"count": &graphql.Field{
					Type: graphql.Int,
				},
			},
		},
	)

	queries = graphql.Fields{
		"vcs": &graphql.Field{
			Type: listType,
			Args: graphql.FieldConfigArgument{
				"team": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Represent the team name or ID",
				},
			},
			Resolve: r.FetchVCSByTeamID,
			// Type: utils.MakeNodeListType("result", utils.MakeListField(vcsType, r.FetchVCSByTeamID)),
		},
		// "vcs": utils.MakeListField(utils.MakeNodeListType("VCSList", vcsType), r.FetchVCSByTeamID),
	}

	mutations = graphql.Fields{}

	return queries, mutations
}
