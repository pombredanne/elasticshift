/*
Copyright 2017 The Elasticshift Authors.
*/
package schema

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/build"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
)

func newBuildSchema(
	ctx context.Context,
	logger logrus.Logger,
	s store.Shift,
) (queries graphql.Fields, mutations graphql.Fields) {

	// return error
	r, _ := build.NewResolver(ctx, logger, s)
	// if err != nil {
	// 	return err
	// }

	buildStatusEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "BuildStatus",
		Values: graphql.EnumValueConfigMap{
			"STUCK": &graphql.EnumValueConfig{
				Value: 1,
			},

			"RUNNING": &graphql.EnumValueConfig{
				Value: 2,
			},

			"SUCCESS": &graphql.EnumValueConfig{
				Value: 3,
			},

			"FAILED": &graphql.EnumValueConfig{
				Value: 4,
			},

			"CANCELLED": &graphql.EnumValueConfig{
				Value: 5,
			},

			"WAITING": &graphql.EnumValueConfig{
				Value: 6,
			},
		},
	})

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Build identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.Build); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"repository_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Repository identifier",
		},

		"vcs_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique identifier for version control system",
		},

		"container_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Container identifier",
		},

		"log": &graphql.Field{
			Type:        graphql.String,
			Description: "Location/path of the build log",
		},

		"started_at": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "Time when the build triggered",
		},

		"ended_at": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "Time when the build completed",
		},

		"triggered_by": &graphql.Field{
			Type:        graphql.String,
			Description: "Show who triggered the build, it could be a pull request (automatially) or user (manually)",
		},

		"status": &graphql.Field{
			Type:        buildStatusEnum,
			Description: "The status of the build",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				if t, ok := p.Source.(types.Build); ok {
					return int(t.Status), nil
				}
				return nil, nil
			},
		},

		"branch": &graphql.Field{
			Type:        graphql.String,
			Description: "The branch to which the build is/was triggered",
		},
	}

	buildType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Build",
			Fields:      fields,
			Description: "An object of Build type",
		},
	)

	buildArgs := graphql.FieldConfigArgument{
		"team": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Build identifier",
		},

		"repository_id": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Status of the build",
		},

		"branch": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "Status of the build",
		},

		"status": &graphql.ArgumentConfig{
			Type:        buildStatusEnum,
			Description: "Status of the build",
		},
	}

	queries = graphql.Fields{
		"build": utils.MakeListType("BuildList", buildType, r.FetchBuild, buildArgs),
	}

	mutations = graphql.Fields{
		"triggerBuild": &graphql.Field{
			Type: buildType,
			Args: graphql.FieldConfigArgument{
				"repository_id": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Name of the system configuration",
				},
				"branch": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "Key for the elasticshift oauth application",
				},
			},
			Resolve: r.TriggerBuild,
		},

		"cancelBuild": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Build identifier",
				},
			},
			Resolve: r.CancelBuild,
		},
	}

	return queries, mutations
}
