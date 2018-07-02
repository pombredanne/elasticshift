/*
Copyright 2018 The Elasticshift Authors.
*/
package store

type infrastructure struct {
	Store
}

// Store provides system level config
type Infrastructure interface {
	Interface
}

// NewStore ..
func newInfrastructureStore(d Database) Infrastructure {
	s := &infrastructure{}
	s.Database = d
	s.CollectionName = "infrastructure"
	return s
}
