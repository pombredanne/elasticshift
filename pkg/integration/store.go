/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	base "gitlab.com/conspico/elasticshift/pkg/store"
	stypes "gitlab.com/conspico/elasticshift/pkg/store/types"
)

type store struct {
	base.Store
}

// Store provides system level config
type Store interface {
	base.Interface
}

// NewStore ..
func NewStore(d stypes.Database) Store {
	s := &store{}
	s.Database = d
	s.CollectionName = "integration"
	return s
}
