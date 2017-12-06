/*
Copyright 2017 The Elasticshift Authors.
*/
package utils

import "github.com/graphql-go/graphql"

func MakeListField(listType graphql.Output, resolve graphql.FieldResolveFn) *graphql.Field {
	return &graphql.Field{
		Type:    listType,
		Resolve: resolve,
		Args: graphql.FieldConfigArgument{
			"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
			"offset": &graphql.ArgumentConfig{Type: graphql.Int},
		},
	}
}

func MakeNodeListType(name string, nodeType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: name,
		Fields: graphql.Fields{
			"nodes": &graphql.Field{Type: graphql.NewList(nodeType)},
			"count": &graphql.Field{Type: graphql.Int},
		},
	})
}
