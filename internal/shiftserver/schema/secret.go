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
	"github.com/elasticshift/elasticshift/internal/shiftserver/secret"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
)

func newSecretSchema(ctx context.Context, loggr logger.Loggr, s store.Shift) (queries graphql.Fields, mutations graphql.Fields) {

	r, _ := secret.NewResolver(ctx, loggr, s)

	secretEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "Kind",
		Values: graphql.EnumValueConfigMap{
			"SECRET": &graphql.EnumValueConfig{
				Value: "secret",
			},

			"SSHKEY": &graphql.EnumValueConfig{
				Value: "sshkey",
			},

			"PGPKEY": &graphql.EnumValueConfig{
				Value: "pgpkey",
			},
		},
	})

	referenceKindEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "ReferenceKind",
		Values: graphql.EnumValueConfigMap{
			"SYS": &graphql.EnumValueConfig{
				Value: "sys",
			},

			"TEAM": &graphql.EnumValueConfig{
				Value: "team",
			},

			"USER": &graphql.EnumValueConfig{
				Value: "user",
			},

			"REPO": &graphql.EnumValueConfig{
				Value: "repo",
			},

			"VCS": &graphql.EnumValueConfig{
				Value: "vcs",
			},
		},
	})

	secretFields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Secret identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.Secret); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the secret",
		},

		"kind": &graphql.Field{
			Type:        secretEnum,
			Description: "Type of secret such as password or key",
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {

			// 	if t, ok := p.Source.(types.Secret); ok {
			// 		return int(t.Kind), nil
			// 	}
			// 	return nil, nil
			// },
		},

		"reference_kind": &graphql.Field{
			Type:        referenceKindEnum,
			Description: "The source to whom the secret belongs to such as team, user, repo etc",
			// Resolve: func(p graphql.ResolveParams) (interface{}, error) {

			// 	if t, ok := p.Source.(types.Secret); ok {
			// 		return int(t.Provider), nil
			// 	}
			// 	return nil, nil
			// },
		},

		"reference_id": &graphql.Field{
			Type:        graphql.String,
			Description: "The identifier belongs to the source reference such as team_id or user_id etc",
		},

		"value": &graphql.Field{
			Type:        graphql.String,
			Description: "The value of the secret",
		},

		"team_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Team Identifier",
		},
	}

	secretType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Secret",
			Fields:      secretFields,
			Description: "An object of Secret type",
		},
	)

	secretArgs := graphql.FieldConfigArgument{

		"team_id": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Team identifier",
		},
	}

	queries = graphql.Fields{
		"Secret": utils.MakeListType("SecretList", secretType, r.FetchSecret, secretArgs),
	}

	mutations = graphql.Fields{

		"addSecret": &graphql.Field{
			Type: secretType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the secret",
				},
				"kind": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(secretEnum),
					Description: "Type of secret such as password or key",
				},
				"reference_kind": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(referenceKindEnum),
					Description: "The source to whom the secret belongs to such as team, user, repo etc",
				},
				"reference_id": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The identifier belongs to the source reference such as team_id or user_id etc",
				},
				"value": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The value of the secret",
				},
				"team_id": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Team identifier",
				},
			},
			Resolve: r.AddSecret,
		},
	}

	return queries, mutations
}
