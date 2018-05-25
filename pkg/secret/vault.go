/*
Copyright 2018 The Elasticshift Authors.
*/
package secret

import (
	"context"
	"fmt"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api/types"
	"gopkg.in/mgo.v2/bson"
)

type Vault interface {
	Encrypt(value string) (string, error)
	Decrypt(value string) (string, error)

	Get(id string) (types.Secret, error)
	Put(secret types.Secret) (string, error)
	Del(id string) error

	GetByReferenceID(id, kind string) (types.Secret, error)
	DelByReferenceID(id, kind string) error
}

type vault struct {
	store  Store
	logger logrus.Logger
	ctx    context.Context
}

func NewVault(s Store, logger logrus.Logger, ctx context.Context) Vault {
	return &vault{
		store:  s,
		logger: logger,
		ctx:    ctx,
	}
}

func (s vault) Get(id string) (types.Secret, error) {

	var sec types.Secret
	err := s.store.FindByID(id, &sec)
	return sec, err
}

func (s vault) Put(sec types.Secret) (string, error) {

	value, err := s.Encrypt(sec.Value)
	if err != nil {
		return "", fmt.Errorf("Error during encryption : %s", err)
	}

	sec.Value = value

	if sec.ID != "" && sec.ID.Hex() != "" {

		// updates
		err = s.store.UpdateId(sec.ID, sec)

	} else {

		// new
		err = s.store.Save(&sec)
	}

	return sec.ID.Hex(), err
}

func (s vault) Del(id string) error {
	return s.store.Remove(id)
}

func (s vault) GetByReferenceID(id, kind string) (types.Secret, error) {

	var sec types.Secret
	err := s.store.FindOne(bson.M{"reference_id": id, "reference_kind": kind}, &sec)
	return sec, err
}

func (s vault) DelByReferenceID(id, kind string) error {
	return s.store.RemoveBySelector(bson.M{"reference_id": id, "reference_kind": kind})
}
