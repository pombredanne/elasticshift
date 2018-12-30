/*
Copyright 2017 The Elasticshift Authors.
*/
package sysconf

import (
	"context"
	"errors"
	"fmt"

	"strings"

	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
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

// Resolver ...
type Resolver interface {
	CreateVCSSysConf(params graphql.ResolveParams) (interface{}, error)
	CreateGenericSysConf(params graphql.ResolveParams) (interface{}, error)
	CreateNFSVolumeSysConf(params graphql.ResolveParams) (interface{}, error)
	FetchVCSSysConf(params graphql.ResolveParams) (interface{}, error)
	FetchVCSSysConfByName(params graphql.ResolveParams) (interface{}, error)
	FetchGenericSysConfByName(params graphql.ResolveParams) (interface{}, error)
	FetchNFSVolumeSysConfByName(params graphql.ResolveParams) (interface{}, error)
	fetchSysconfByName(kind string, params graphql.ResolveParams, result interface{}) error
}

type resolver struct {
	store  store.Sysconf
	logger *logrus.Entry
}

// NewResolver ...
func NewResolver(ctx context.Context, loggr logger.Loggr, s store.Shift) (Resolver, error) {

	r := &resolver{
		store:  s.Sysconf,
		logger: loggr.GetLogger("graphql/sysconf"),
	}
	return r, nil
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
		return nil, fmt.Errorf("Failed to create vcs sysconf: %v", err)
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
		return nil, fmt.Errorf("Failed to create generic sysconf: %v", err)
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
		return nil, fmt.Errorf("Failed to create NFS volume sysconf: %v", err)
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
