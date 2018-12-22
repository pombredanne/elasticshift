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
	"github.com/elasticshift/elasticshift/internal/shiftserver/container"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
)

func newContainerSchema(ctx context.Context, loggr logger.Loggr, s store.Shift) (queries graphql.Fields, mutations graphql.Fields) {

	// return error
	r, _ := container.NewResolver(ctx, loggr, s)
	// if err != nil {
	// 	return err
	// }

	containerStatusEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "ContainerStatus",
		Values: graphql.EnumValueConfigMap{
			"RUNNING": &graphql.EnumValueConfig{
				Value: 1,
			},

			"STARTED": &graphql.EnumValueConfig{
				Value: 2,
			},

			"STOPPED": &graphql.EnumValueConfig{
				Value: 3,
			},
		},
	})

	fields := graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.ID,
			Description: "Container identifier",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if t, ok := p.Source.(types.Container); ok {
					return t.ID.Hex(), nil
				}
				return nil, nil
			},
		},

		"build_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Build identifier",
		},

		"repository_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Repository Identifier",
		},

		"container_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Container identifier",
		},

		"vcs_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Version Control system Identifier",
		},

		"orchestration_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Orchestration Identifier",
		},

		"image": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "Name of the contianer image",
		},

		"started_at": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "Time when the container started",
		},

		"stopped_at": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "Time when the container stopped",
		},

		"duration": &graphql.Field{
			Type:        graphql.String,
			Description: "The lifetime of the container",
		},

		"kind": &graphql.Field{
			Type:        graphql.String,
			Description: "The lifetime of the container",
		},

		"status": &graphql.Field{
			Type:        containerStatusEnum,
			Description: "The status of the container",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {

				if t, ok := p.Source.(types.Container); ok {
					return int(t.Status), nil
				}
				return nil, nil
			},
		},
	}

	containerType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:        "Container",
			Fields:      fields,
			Description: "An object of Container type",
		},
	)

	containerArgs := graphql.FieldConfigArgument{

		"id": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Container identifier",
		},

		"team": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Team identifier",
		},

		"build_id": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "Build identifier",
		},

		"status": &graphql.ArgumentConfig{
			Type:        containerStatusEnum,
			Description: "Status of the container",
		},
	}

	queries = graphql.Fields{
		"container": utils.MakeListType("ContainerList", containerType, r.FetchContainer, containerArgs),
	}

	mutations = graphql.Fields{}

	return queries, mutations
}
