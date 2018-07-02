/*
Copyright 2018 The Elasticshift Authors.
*/
package store

type integration struct {
	Store
}

// Store provides system level config
type Integration interface {
	Interface
}

// NewStore ..
func newIntegrationStore(d Database) Integration {
	s := &integration{}
	s.Database = d
	s.CollectionName = "integration"
	return s
}
