/*
Copyright 2018 The Elasticshift Authors.
*/
package store

type plugin struct {
	Store
}

// Store provides system level config
type Plugin interface {
	Interface
}

// NewStore ..
func newPluginStore(d Database) Plugin {
	s := &plugin{}
	s.Database = d
	s.CollectionName = "plugin"
	return s
}
