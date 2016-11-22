// Package core ...
// Author: Ghazni Nattarshah
// Date: Nov 22, 2016
package core

import (
	"strings"

	"github.com/jinzhu/gorm"
	"gitlab.com/conspico/esh"
)

var (
	recordNotFound = "record not found"
)

// Datastore ..
// A base datasource that performs actualy sql interactions.
type datastore struct {
	db        *gorm.DB
	operation chan func(db *gorm.DB)
}

const (
	ok = true
)

// Create ..
func (ds *datastore) Create(model interface{}) error {

	result := make(chan interface{}, 1)
	ds.operation <- func(db *gorm.DB) {

		db.NewRecord(model)
		e := db.Create(model).Error
		if e != nil {
			result <- e
		}
		result <- ok
	}

	return grabResult(<-result)
}

// Update ..
func (ds *datastore) Update(old interface{}, updated interface{}) error {

	result := make(chan interface{}, 1)
	ds.operation <- func(db *gorm.DB) {

		e := db.Model(old).Updates(updated).Error
		if e != nil {
			result <- e
		}
		result <- ok
	}
	return grabResult(<-result)
}

// Read ..
func (ds *datastore) Read(query string, value interface{}, params ...interface{}) error {

	result := make(chan interface{}, 1)
	ds.operation <- func(db *gorm.DB) {

		e := db.Raw(query, params...).Scan(value).Error
		if e != nil && !strings.EqualFold(recordNotFound, e.Error()) {
			result <- e
		}
		result <- ok
	}
	return grabResult(<-result)
}

// Delete ..
func (ds *datastore) Delete(model interface{}) error {

	result := make(chan interface{}, 1)
	ds.operation <- func(db *gorm.DB) {

		e := db.Delete(model).Error
		if e != nil {
			result <- e
		}
		result <- ok
	}
	return grabResult(<-result)
}

// DeleteMultiple ..
func (ds *datastore) DeleteMultiple(model interface{}, ids []string) error {

	result := make(chan interface{}, 1)
	ds.operation <- func(db *gorm.DB) {

		e := db.Delete(model, "ID IN (?)", ids).Error
		if e != nil {
			result <- e
		}
		result <- ok
	}
	return grabResult(<-result)
}

// Process ..
func (ds *datastore) process() {

	for op := range ds.operation {
		op(ds.db)
	}
}

func grabResult(value interface{}) error {

	switch v := value.(type) {
	case error:
		return v
	case bool:
		return nil
	default:
		return nil
	}
}

// NewDatasource ..
// Create a new base datasource
func NewDatasource(db *gorm.DB) esh.Datastore {

	ds := &datastore{db: db, operation: make(chan func(db *gorm.DB))}
	go ds.process()
	return ds
}
