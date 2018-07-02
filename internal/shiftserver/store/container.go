/*
Copyright 2018 The Elasticshift Authors.
*/
package store

type container struct {
	Store
}

// Container ...
// Store provides container related operation
type Container interface {
	Interface
}

// NewStore ..
func newContainerStore(d Database) Container {
	s := &container{}
	s.Database = d
	s.CollectionName = "container"
	return s
}
