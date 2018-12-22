/*
Copyright 2018 The Elasticshift Authors.
*/
package schema

import (
	"context"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	"github.com/elasticshift/elasticshift/internal/shiftserver/identity/oauth2/providers"
	"github.com/elasticshift/elasticshift/internal/shiftserver/pubsub"
	"github.com/elasticshift/elasticshift/internal/shiftserver/resolver"
	"github.com/elasticshift/elasticshift/internal/shiftserver/secret"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
)

// Construct ...
// Construct the graphql schema
func Construct(
	ctx context.Context,
	loggr logger.Loggr,
	providers providers.Providers,
	s store.Shift,
	vault secret.Vault,
	pb pubsub.Engine,
	r *resolver.Shift,
) (*graphql.Schema, error) {

	// initialize schema
	queries := graphql.Fields{}
	mutations := graphql.Fields{}
	subscriptions := graphql.Fields{}

	// team fields
	teamQ, teamM := newTeamSchema(ctx, loggr, s)
	appendFields(queries, teamQ)
	appendFields(mutations, teamM)

	// vcs fields
	vcsQ, vcsM := newVcsSchema(ctx, loggr, providers, s)
	appendFields(queries, vcsQ)
	appendFields(mutations, vcsM)

	// repository fields
	repositoryQ, repositoryM := newRepositorySchema(ctx, loggr, s)
	appendFields(queries, repositoryQ)
	appendFields(mutations, repositoryM)

	// sysconf fields
	sysconfQ, sysconfM := newSysconfSchema(ctx, loggr, s)
	appendFields(queries, sysconfQ)
	appendFields(mutations, sysconfM)

	// build fields
	buildQ, buildM, buildS := newBuildSchema(ctx, loggr, s, pb, r)
	appendFields(queries, buildQ)
	appendFields(mutations, buildM)
	appendFields(subscriptions, buildS)

	// app fields
	pluginQ, pluginM := newPluginSchema(ctx, loggr, s)
	appendFields(queries, pluginQ)
	appendFields(mutations, pluginM)

	// container fields
	containerQ, containerM := newContainerSchema(ctx, loggr, s)
	appendFields(queries, containerQ)
	appendFields(mutations, containerM)

	// integration fields
	integrationQ, integrationM := newIntegrationSchema(ctx, loggr, s)
	appendFields(queries, integrationQ)
	appendFields(mutations, integrationM)

	// infrastructure fields
	infrastructureQ, infrastructureM := newInfrastructureSchema(ctx, loggr, s)
	appendFields(queries, infrastructureQ)
	appendFields(mutations, infrastructureM)

	// default fields
	defaultQ, defaultM := newDefaultsSchema(ctx, loggr, s)
	appendFields(queries, defaultQ)
	appendFields(mutations, defaultM)

	// secret fields
	secretQ, secretM := newSecretSchema(ctx, loggr, s)
	appendFields(queries, secretQ)
	appendFields(mutations, secretM)

	// secret fields
	shiftfileQ, shiftfileM := newShiftfileSchema(ctx, loggr, s)
	appendFields(queries, shiftfileQ)
	appendFields(mutations, shiftfileM)

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queries}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutations}
	rootSubscription := graphql.ObjectConfig{Name: "RootSubscription", Fields: subscriptions}

	schemaConfig := graphql.SchemaConfig{
		Query:        graphql.NewObject(rootQuery),
		Mutation:     graphql.NewObject(rootMutation),
		Subscription: graphql.NewObject(rootSubscription),
	}

	schm, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create schema due to errors %v", err)
	}
	return &schm, nil
}

// Utility method to append fields
func appendFields(fields graphql.Fields, input graphql.Fields) {

	for k, v := range input {
		fields[k] = v
	}
}
