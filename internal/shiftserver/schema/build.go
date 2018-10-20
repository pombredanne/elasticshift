/*
Copyright 2017 The Elasticshift Authors.
*/
package schema

import (
	"context"

	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/build"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/pubsub"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/resolver"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
)

// type & fields
var (
	buildStatusEnum = graphql.NewEnum(graphql.EnumConfig{
		Name: "BuildStatus",
		Values: graphql.EnumValueConfigMap{

			"WAITING": &graphql.EnumValueConfig{
				Value: 1,
			},

			"PREPARING": &graphql.EnumValueConfig{
				Value: 2,
			},

			"RUNNING": &graphql.EnumValueConfig{
				Value: 3,
			},

			"SUCCESS": &graphql.EnumValueConfig{
				Value: 4,
			},

			"FAILED": &graphql.EnumValueConfig{
				Value: 5,
			},

			"CANCELLED": &graphql.EnumValueConfig{
				Value: 6,
			},

			"STUCK": &graphql.EnumValueConfig{
				Value: 7,
			},
		},
	})

	fields = graphql.Fields{
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

					var status int
					switch t.Status {
					case types.BuildStatusWaiting: // waiting
						status = 1
					case types.BuildStatusPreparing: // preparing
						status = 2
					case types.BuildStatusRunning: // running
						status = 3
					case types.BuildStatusSuccess: // success
						status = 4
					case types.BuildStatusFailed: // failed
						status = 5
					case types.BuildStatusCancel: // cancelled
						status = 6
					case types.BuildStatusStuck: // stuck
						status = 7
					}

					return status, nil
				}
				return nil, nil
			},
		},

		"branch": &graphql.Field{
			Type:        graphql.String,
			Description: "The branch to which the build is/was triggered",
		},

		"graph": &graphql.Field{
			Type:        graphql.String,
			Description: "Flow graph",
		},

		"reason": &graphql.Field{
			Type:        graphql.String,
			Description: "Reason for build failures.",
		},

		"duration": &graphql.Field{
			Type:        graphql.String,
			Description: "Duration of the actual build time",
		},
	}

	BuildType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Build",
			Fields:      fields,
			Description: "An object of Build type",
		},
	)
)

func newBuildSchema(
	ctx context.Context,
	loggr logger.Loggr,
	s store.Shift,
	ps pubsub.Engine,
	rs *resolver.Shift,
) (queries graphql.Fields, mutations graphql.Fields, subscriptions graphql.Fields) {

	r, _ := build.NewResolver(ctx, loggr, s, ps)
	rs.Build = r

	buildArgs := graphql.FieldConfigArgument{
		"team": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Build identifier",
		},

		"repository_id": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Status of the build",
		},

		"id": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "Build Identifier",
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
		"build": utils.MakeListType("BuildList", BuildType, r.FetchBuild, buildArgs),
	}

	mutations = graphql.Fields{
		"triggerBuild": &graphql.Field{
			Type: BuildType,
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

	subscriptions = graphql.Fields{
		pubsub.SubscribeBuildUpdate: &graphql.Field{
			Type: BuildType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Build Identifier",
				},
			},
			Resolve: r.FetchBuildByID,
		},
	}

	return queries, mutations, subscriptions
}
