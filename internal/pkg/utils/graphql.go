/*
Copyright 2017 The Elasticshift Authors.
*/
package utils

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func MakeListType(name string, listType graphql.Output, resolve graphql.FieldResolveFn, args graphql.FieldConfigArgument) *graphql.Field {

	obj := graphql.NewObject(graphql.ObjectConfig{
		Name: name,
		Fields: graphql.Fields{
			"nodes": &graphql.Field{Type: graphql.NewList(listType)},
			"count": &graphql.Field{Type: graphql.Int},
		},
	})

	field := &graphql.Field{
		Type:    obj,
		Resolve: resolve,
		Args: graphql.FieldConfigArgument{
			"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
			"offset": &graphql.ArgumentConfig{Type: graphql.Int},
		},
	}

	// Append the additional query param
	for k, v := range args {
		field.Args[k] = v
	}

	return field
}

// Errors convert from GraphQL errors to regular errors.
func GraphQLErrors(errors []gqlerrors.FormattedError) []error {
	if len(errors) == 0 {
		return nil
	}

	out := make([]error, len(errors))
	for i := range errors {
		out[i] = errors[i]
	}
	return out
}
