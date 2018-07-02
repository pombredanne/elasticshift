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
	"gitlab.com/conspico/elasticshift/internal/store"
)

var (
	errNameIsRequired        = errors.New("Can't fetch setting with out a name.")
	errVCSAlreadyExist       = errors.New("VCSSysConf name already exist")
	errGenericAlreadyExist   = errors.New("GenericSysConf name already exist")
	errNFSVolumeAlreadyExist = errors.New("NFSVolumeSysConf name already exist")
)

const (
	VcsKind       = "vcs"
	GenericKind   = "generic"
	VolumeNfsKind = "nfs"

	DEFAULT_STORAGE = "default-storage"
)

type resolver struct {
	store  store.Sysconf
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
	res.Kind = VcsKind

	result, err := r.FetchVCSSysConfByName(params)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return nil, fmt.Errorf("Failed to create vcs sysconf", err)
	}

	if result.(types.VCSSysConf).Name != "" {
		return nil, errVCSAlreadyExist
	}

	err = r.store.Save(res)
	return res, err
}

func (r *resolver) CreateGenericSysConf(params graphql.ResolveParams) (interface{}, error) {

	name, _ := params.Args["name"].(string)
	value, _ := params.Args["value"].(string)

	res := &types.GenericSysConf{}
	res.Name = name
	res.Value = value
	res.Kind = GenericKind

	result, err := r.FetchGenericSysConfByName(params)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return nil, fmt.Errorf("Failed to create generic sysconf", err)
	}

	if result.(types.GenericSysConf).Name != "" {
		return nil, errGenericAlreadyExist
	}

	err = r.store.Save(res)
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
	res.Kind = VolumeNfsKind

	result, err := r.FetchNFSVolumeSysConfByName(params)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return nil, fmt.Errorf("Failed to create NFS volume sysconf", err)
	}

	if result.(types.NFSVolumeSysConf).Name != "" {
		return nil, errNFSVolumeAlreadyExist
	}

	err = r.store.Save(res)
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
	err := r.fetchSysconfByName(VcsKind, params, &result)
	return result, err
}

func (r *resolver) FetchGenericSysConfByName(params graphql.ResolveParams) (interface{}, error) {
	var result types.GenericSysConf
	err := r.fetchSysconfByName(GenericKind, params, &result)
	return result, err
}

func (r *resolver) FetchNFSVolumeSysConfByName(params graphql.ResolveParams) (interface{}, error) {
	var result types.NFSVolumeSysConf
	err := r.fetchSysconfByName(VolumeNfsKind, params, &result)
	return result, err
}

func (r *resolver) fetchSysconfByName(kind string, params graphql.ResolveParams, result interface{}) error {

	name, _ := params.Args["name"].(string)
	if name == "" {
		return errNameIsRequired
	}

	return r.store.GetSysConf(kind, name, result)
}
