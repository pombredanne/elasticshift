/*
Copyright 2018 The Elasticshift Authors.
*/
package defaults

import (
	"context"
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/store"
)

func InitSchema(logger logrus.Logger, ctx context.Context, s store.Shift) (queries graphql.Fields, mutations graphql.Fields) {

	r := &resolver{
		store:     s.Defaults,
		logger:    logger,
		Ctx:       ctx,
		teamStore: s.Team,
	}

	// propertyFields := graphql.Fields{
	// 	"key": &graphql.Field{
	// 		Type:        graphql.String,
	// 		Description: "Property key",
	// 	},
	// 	"value": &graphql.Field{
	// 		Type:        graphql.String,
	// 		Description: "Property value",
	// 	},
	// }

	// propertyType := graphql.NewObject(
	// 	graphql.ObjectConfig{
	// 		Name:        "Property",
	// 		Fields:      propertyFields,
	// 		Description: "An object of Property type",
	// 	},
	// )

	kindEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "DefaultKind",
		Values: graphql.EnumValueConfigMap{
			"TEAM": &graphql.EnumValueConfig{
				Value: 1,
			},

			"USER": &graphql.EnumValueConfig{
				Value: 2,
			},
		},
	})

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Default identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(*types.Default); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"reference_id": &graphql.Field{
			Type:        graphql.String,
			Description: "The reference identifier which is the id of team or user",
		},

		"kind": &graphql.Field{
			Type:        kindEnum,
			Description: "The kind for which the defaults belongs to such as user or team",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				if t, ok := p.Source.(*types.Default); ok {
					return int(t.Kind), nil
				}
				return nil, nil
			},
		},

		"container_engine_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Container Engine identifier",
		},

		"storage_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Storage identifier",
		},

		"languages": &graphql.Field{
			Type:        graphql.String,
			Description: "Default Language specification (shiftfile)",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				if t, ok := p.Source.(*types.Default); ok {

					data, err := json.Marshal(t.Languages)
					if err != nil {
						return nil, err
					}
					return string(data), nil
				}
				return nil, nil
			},
		},
	}

	defaultType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Default",
			Fields:      fields,
			Description: "An object of Default type",
		},
	)

	defaultArgs := graphql.FieldConfigArgument{

		"team": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "Team identifier",
		},

		"id": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "Default identifier",
		},

		"kind": &graphql.ArgumentConfig{
			Type:        kindEnum,
			Description: "The kind for which the defaults belongs to such as user or team",
		},

		"reference_id": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Build identifier",
		},
	}

	queries = graphql.Fields{

		"default": &graphql.Field{
			Type:    defaultType,
			Args:    defaultArgs,
			Resolve: r.FetchDefault,
		},
	}

	mutations = graphql.Fields{

		"setDefaults": &graphql.Field{
			Type: defaultType,
			Args: graphql.FieldConfigArgument{
				"reference_id": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Reference Identifier could be team or user id",
				},
				"kind": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(kindEnum),
					Description: "The kind that defaults belongs to such as user or team",
				},
				"container_engine_id": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Default container engine identifier",
				},
				"storage_id": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Default storage identifier",
				},
				"languages": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Default language description file (shift file), this should be in json format with key value pair",
				},
			},
			Resolve: r.SetDefaults,
		},
	}

	return queries, mutations
}
