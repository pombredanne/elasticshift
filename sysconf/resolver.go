/*
Copyright 2017 The Elasticshift Authors.
*/
package sysconf

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
)

var (
	errNameIsRequired = errors.New("Can't fetch setting with out a name.")
)

type resolver struct {
	store  Store
	logger logrus.Logger
}

func (r *resolver) CreateVCSSysConf(params graphql.ResolveParams) (interface{}, error) {

	name, _ := params.Args["name"].(string)
	key, _ := params.Args["key"].(string)
	secret, _ := params.Args["secret"].(string)
	callbackURL, _ := params.Args["callbackURL"].(string)

	res := &types.VCSSysConf{}
	res.Name = name
	res.Key = key
	res.Secret = secret
	res.CallbackURL = callbackURL

	err := r.store.SaveVCSSysConf(res)

	return res, err
}

func (r *resolver) FetchVCSSysConfByName(params graphql.ResolveParams) (interface{}, error) {

	name, _ := params.Args["name"].(string)
	if name == "" {
		return nil, errNameIsRequired
	}

	res, err := r.store.GetVCSSysConfByName(name)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *resolver) FetchVCSSysConf(params graphql.ResolveParams) (interface{}, error) {
	result, err := r.store.GetVCSSysConf()
	return &result, err
}
