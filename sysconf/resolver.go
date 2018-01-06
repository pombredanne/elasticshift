/*
Copyright 2017 The Elasticshift Authors.
*/
package sysconf

import (
	"errors"
	"fmt"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
)

var (
	errNameIsRequired        = errors.New("Can't fetch setting with out a name.")
	errVCSAlreadyExist       = errors.New("VCSSysConf name already exist")
	errGenericAlreadyExist   = errors.New("GenericSysConf name already exist")
	errNFSVolumeAlreadyExist = errors.New("NFSVolumeSysConf name already exist")
)

const (
	vcsKind       = "vcs"
	genericKind   = "generic"
	volumeNfsKind = "nfs"
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
	res.Kind = vcsKind

	result, err := r.FetchVCSSysConfByName(params)
	if err != nil {
		return nil, fmt.Errorf("Failed to create vcs sysconf", err)
	}

	if result != nil {
		return nil, errVCSAlreadyExist
	}

	err = r.store.SaveSysConf(res)
	return res, err
}

func (r *resolver) CreateGenericSysConf(params graphql.ResolveParams) (interface{}, error) {

	name, _ := params.Args["name"].(string)
	value, _ := params.Args["value"].(string)

	res := &types.GenericSysConf{}
	res.Name = name
	res.Value = value
	res.Kind = genericKind

	result, err := r.FetchGenericSysConfByName(params)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return nil, fmt.Errorf("Failed to create generic sysconf", err)
	}

	if result.(types.GenericSysConf).Name != "" {
		return nil, errGenericAlreadyExist
	}

	err = r.store.SaveSysConf(res)
	return res, err

}

func (r *resolver) CreateNFSVolumeSysConf(params graphql.ResolveParams) (interface{}, error) {

	name, _ := params.Args["name"].(string)
	server, _ := params.Args["server"].(string)
	accessMode, _ := params.Args["accessmode"].(int)

	res := &types.NFSVolumeSysConf{}
	res.Name = name
	res.Server = server
	res.AccessMode = accessMode
	res.Kind = volumeNfsKind

	result, err := r.FetchNFSVolumeSysConfByName(params)
	if err != nil {
		return nil, fmt.Errorf("Failed to create NFS volume sysconf", err)
	}

	if result != nil {
		return nil, errNFSVolumeAlreadyExist
	}

	err = r.store.SaveSysConf(res)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (r *resolver) FetchVCSSysConf(params graphql.ResolveParams) (interface{}, error) {
	result, err := r.store.GetVCSSysConf()
	return result, err
}

func (r *resolver) FetchVCSSysConfByName(params graphql.ResolveParams) (interface{}, error) {
	var result types.VCSSysConf
	err := r.fetchSysconfByName(vcsKind, params, &result)
	return result, err
}

func (r *resolver) FetchGenericSysConfByName(params graphql.ResolveParams) (interface{}, error) {
	var result types.GenericSysConf
	err := r.fetchSysconfByName(genericKind, params, &result)
	return result, err
}

func (r *resolver) FetchNFSVolumeSysConfByName(params graphql.ResolveParams) (interface{}, error) {
	var result types.NFSVolumeSysConf
	err := r.fetchSysconfByName(volumeNfsKind, params, &result)
	return result, err
}

func (r *resolver) fetchSysconfByName(kind string, params graphql.ResolveParams, result interface{}) error {

	name, _ := params.Args["name"].(string)
	if name == "" {
		return errNameIsRequired
	}

	return r.store.GetSysConf(kind, name, result)
}
