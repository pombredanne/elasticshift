/*
Copyright 2018 The Elasticshift Authors.
*/
package schema

import (
	"context"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/internal/pkg/identity/oauth2/providers"
	"gitlab.com/conspico/elasticshift/internal/pkg/secret"
	"gitlab.com/conspico/elasticshift/internal/store"
)

// Construct ...
// Construct the graphql schema
func Construct(
	ctx context.Context,
	logger logrus.Logger,
	providers providers.Providers,
	s store.Shift,
	vault secret.Vault,
) (*graphql.Schema, error) {

	// initialize schema
	queries := graphql.Fields{}
	mutations := graphql.Fields{}

	// team fields
	teamQ, teamM := newTeamSchema(ctx, logger, s)
	appendFields(queries, teamQ)
	appendFields(mutations, teamM)

	// vcs fields
	vcsQ, vcsM := newVcsSchema(ctx, logger, providers, s)
	appendFields(queries, vcsQ)
	appendFields(mutations, vcsM)

	// repository fields
	repositoryQ, repositoryM := newRepositorySchema(ctx, logger, s)
	appendFields(queries, repositoryQ)
	appendFields(mutations, repositoryM)

	// sysconf fields
	sysconfQ, sysconfM := newSysconfSchema(ctx, logger, s)
	appendFields(queries, sysconfQ)
	appendFields(mutations, sysconfM)

	// build fields
	buildQ, buildM := newBuildSchema(ctx, logger, s)
	appendFields(queries, buildQ)
	appendFields(mutations, buildM)

	// app fields
	pluginQ, pluginM := newPluginSchema(ctx, logger, s)
	appendFields(queries, pluginQ)
	appendFields(mutations, pluginM)

	// container fields
	containerQ, containerM := newContainerSchema(ctx, logger, s)
	appendFields(queries, containerQ)
	appendFields(mutations, containerM)

	// integration fields
	integrationQ, integrationM := newIntegrationSchema(ctx, logger, s)
	appendFields(queries, integrationQ)
	appendFields(mutations, integrationM)

	// infrastructure fields
	infrastructureQ, infrastructureM := newInfrastructureSchema(ctx, logger, s)
	appendFields(queries, infrastructureQ)
	appendFields(mutations, infrastructureM)

	// default fields
	defaultQ, defaultM := newDefaultsSchema(ctx, logger, s)
	appendFields(queries, defaultQ)
	appendFields(mutations, defaultM)

	// secret fields
	secretQ, secretM := newSecretSchema(ctx, logger, s)
	appendFields(queries, secretQ)
	appendFields(mutations, secretM)

	// secret fields
	shiftfileQ, shiftfileM := newShiftfileSchema(ctx, logger, s)
	appendFields(queries, shiftfileQ)
	appendFields(mutations, shiftfileM)

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queries}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutations}

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: graphql.NewObject(rootMutation),
	}

	schm, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create schema due to errors :v", err)
	}
	return &schm, nil
}

// Utility method to append fields
func appendFields(fields graphql.Fields, input graphql.Fields) {

	for k, v := range input {
		fields[k] = v
	}
}
