/*
Copyright 2018 The Elasticshift Authors.
*/
package store

type shiftfile struct {
	Store
}

// Store provides system level config
type Shiftfile interface {
	Interface
}

// NewStore ..
func newShiftfileStore(d Database) Shiftfile {
	s := &shiftfile{}
	s.Database = d
	s.CollectionName = "shiftfile"
	return s
}
