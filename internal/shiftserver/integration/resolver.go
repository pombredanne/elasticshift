/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	"github.com/elasticshift/elasticshift/internal/shiftserver/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errIDCantBeEmpty     = errors.New("Integration ID cannot be empty")
	errTeamCannotBeEmpty = errors.New("Team must be provided")
)

const (
	INT_ContainerEngine int = iota + 1
	INT_Storage
)

// Resolver ...
type Resolver interface {
	FetchContainerEngine(params graphql.ResolveParams) (interface{}, error)
	FetchStorage(params graphql.ResolveParams) (interface{}, error)
	AddContainerEngine(params graphql.ResolveParams) (interface{}, error)
	AddStorage(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	store  store.Integration
	logger *logrus.Entry
	Ctx    context.Context
	sysconfStore store.Sysconf
}

// NewResolver ...
func NewResolver(ctx context.Context, loggr logger.Loggr, s store.Shift) (Resolver, error) {

	r := &resolver{
		store:  s.Integration,
		logger: loggr.GetLogger("graphql/integration"),
		Ctx:    ctx,
		sysconfStore: s.Sysconf,
	}
	return r, nil
}

func (r *resolver) FetchContainerEngine(params graphql.ResolveParams) (interface{}, error) {

	team, _ := params.Args["team"].(string)
	if team == "" {
		return nil, errTeamCannotBeEmpty
	}

	q := bson.M{"team": team, "internal_type": INT_ContainerEngine}

	id, _ := params.Args["id"].(string)
	if id != "" {
		q["_id"] = bson.ObjectIdHex(id)
	}

	var err error
	var result []types.ContainerEngine
	r.store.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	var res types.ContainerEngineList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}

func (r *resolver) FetchStorage(params graphql.ResolveParams) (interface{}, error) {

	team, _ := params.Args["team"].(string)
	if team == "" {
		return nil, errTeamCannotBeEmpty
	}

	q := bson.M{"team": team, "internal_type": INT_Storage}

	id, _ := params.Args["id"].(string)
	if id != "" {
		q["_id"] = bson.ObjectIdHex(id)
	}

	var err error
	var result []types.Storage
	r.store.Execute(func(c *mgo.Collection) {
		err = c.Find(q).All(&result)
	})

	var res types.StorageList
	res.Nodes = result
	res.Count = len(res.Nodes)

	return &res, err
}

func (r *resolver) AddContainerEngine(params graphql.ResolveParams) (interface{}, error) {

	team, _ := params.Args["team"].(string)
	name, _ := params.Args["name"].(string)

	var ce types.ContainerEngine
	err := r.store.FindOne(bson.M{"team": team, "name": name}, &ce)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return nil, fmt.Errorf("Failed to check if the given integration already exist :%v", err)
	}

	if ce.ID.Hex() != "" {
		return nil, fmt.Errorf("The container engine name '%s' already exist for your team", name)
	}

	kind, _ := params.Args["kind"].(int)
	provider, _ := params.Args["provider"].(int)
	host, _ := params.Args["host"].(string)
	certificate, _ := params.Args["certificate"].(string)
	token, _ := params.Args["token"].(string)
	version, _ := params.Args["version"].(string)

	i := types.ContainerEngine{}
	i.Name = name
	i.Team = team
	i.Kind = kind
	i.Host = host
	i.Certificate = certificate
	i.Token = token
	i.Provider = provider
	i.InternalType = INT_ContainerEngine
	i.Version = version

	err = r.store.Save(&i)
	if err != nil {
		return nil, fmt.Errorf("Failed to add integration: %v", err)
	}
	return i, nil
}

func (r *resolver) AddStorage(params graphql.ResolveParams) (interface{}, error) {

	args := params.Args["storage"].(map[string]interface{})

	name, _ := args["name"].(string)
	team, _ := args["team"].(string)

	var stor types.Storage
	err := r.store.FindOne(bson.M{"team": team, "name": name}, &stor)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return nil, fmt.Errorf("Failed to check if the given storage integration already exist :%v", err)
	}

	if stor.ID.Hex() != "" {
		return nil, fmt.Errorf("The storage name '%s' already exist for your team", name)
	}

	kind, _ := args["kind"].(int)
	provider, _ := args["provider"].(int)

	source, _ := args["storage_source"].(map[string]interface{})
	if len(source) == 0 {
		return nil, fmt.Errorf("Storage source must be provided")
	}

	i := types.Storage{}
	i.ID = bson.NewObjectId()
	i.Name = name
	i.Team = team
	i.Kind = kind
	i.Provider = provider

	sourceType := types.StorageSource{}
	if val, ok := source["nfs"].(map[string]interface{}); ok {

		nfsType := &types.NFSStorage{}
		nfsType.Server, _ = val["server"].(string)
		nfsType.Path, _ = val["path"].(string)
		nfsType.MountPath, _ = val["mount_path"].(string)
		nfsType.ReadOnly, _ = val["readonly"].(bool)

		sourceType.NFS = nfsType
	} else if val, ok := source["minio"].(map[string]interface{}); ok {

		minioType := &types.MinioStorage{}
		minioType.Host, _ = val["host"].(string)
		minioType.Certificate, _ = val["certificate"].(string)
		minioType.AccessKey, _ = val["accesskey"].(string)
		minioType.SecretKey, _ = val["secretkey"].(string)
		minioType.BucketName, _ = val["bucket_name"].(string)

		sourceType.Minio = minioType
	}
	i.StorageSource = sourceType
	i.InternalType = INT_Storage

	var storag StorageInterface
	if kind != NFS {

		// setup the storage.
		storag, err = NewStorage(r.logger, i)
		if err != nil {
			return nil, fmt.Errorf("Failed to connect : %v", err)
		}
	}

	err = r.store.Save(&i)
	if err != nil {
		return nil, fmt.Errorf("Failed to add integration: %v", err)
	}

	// background goroutine to setup storage.
	go r.setupStorage(storag, i)

	return i, nil
}

func (r *resolver) setupStorage(stor StorageInterface, s types.Storage) {

	workerURL := r.getWorkerURL()

	// TODO get the bucketname as part of storage
	var bucketName string 
	if s.Minio.BucketName != "" {
		bucketName = s.Minio.BucketName
	} else if s.Name != "" {
		bucketName = s.Name
	} else {
		bucketName = "elasticshift"
	}

	objectName, err := stor.SetupStorage(bucketName, workerURL)
	if err != nil {
		s.Reason = err.Error()
		err = r.store.Update(s.ID, &s)
		if err != nil {
			r.logger.Errorf("Failed to update the reason while setting up the storage. %v", err)
		}
		return
	}

	err = r.store.UpdateWorkerPath(s.ID, objectName)
	if err != nil {
		r.logger.Errorf("Failed to update worker path: %v", err)
	}
}

func (r *resolver) getWorkerURL() string {

	// TODO fetch the url from sysconf
	var result types.GenericSysConf
	err := r.sysconfStore.GetSysConf(store.GenericKind, "worker.url", &result)
	if err != nil {
		return "" 
	}
	
	return result.Value	
}
