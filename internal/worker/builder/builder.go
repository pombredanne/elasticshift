/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"context"
	"fmt"
	"io"
	"time"

	"path/filepath"

	"github.com/pkg/errors"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/pkg/graph"
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/ast"
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/parser"
	"gitlab.com/conspico/elasticshift/internal/pkg/storage"
	"gitlab.com/conspico/elasticshift/internal/pkg/vcs"
	"gitlab.com/conspico/elasticshift/internal/worker/logshipper"
	wtypes "gitlab.com/conspico/elasticshift/internal/worker/types"
	"google.golang.org/grpc"
)

var (
	DIR_CODE    = "code"
	DIR_PLUGINS = "plugins"
	DIR_WORKER  = "worker"
	DIR_LOGS    = "logs"

	// TODO check for windows container
	VOL_SHIFT   = "/opt/elasticshift"
	VOL_CODE    = filepath.Join(VOL_SHIFT, DIR_CODE)
	VOL_PLUGINS = filepath.Join(VOL_SHIFT, DIR_PLUGINS)
	VOL_LOGS    = filepath.Join(VOL_SHIFT, DIR_LOGS)
)

type builder struct {
	shiftconn   *grpc.ClientConn
	ctx         context.Context
	wctx        wtypes.Context
	config      wtypes.Config
	shiftclient api.ShiftClient
	project     *api.GetProjectRes
	storage     *storage.ShiftStorage
	logshipper  logshipper.LogShipper

	f *ast.File
	g *graph.Graph

	done chan int

	writer io.Writer
}

func New(ctx wtypes.Context, shiftconn *grpc.ClientConn, writer io.Writer, done chan int) error {

	b := builder{}
	b.shiftconn = shiftconn
	b.ctx = ctx.Context
	b.wctx = ctx
	b.shiftclient = ctx.Client
	b.config = ctx.Config
	b.done = done
	b.writer = writer

	return b.run()
}

func (b *builder) run() error {

	// Get the project information
	proj, err := b.shiftclient.GetProject(b.ctx, &api.GetProjectReq{BuildId: b.config.BuildID, IncludeShiftfile: !b.config.RepoBasedShiftFile})
	if err != nil {
		return fmt.Errorf("Failed to get the project/repository detail from shift server: %v\n", err)
	}
	b.project = proj

	b.wctx.EnvLogger.Printf("Project Info: %v\n", proj)

	// 1. Ensure connection to log storage is good, this container should be loaded with

	// 2. Load the build cache, if available ensure it

	// 3. Fetch the shiftfile

	// 4. otherwise use the global language spec defined by elasticshift
	var f []byte
	if b.config.RepoBasedShiftFile {
		f, err = vcs.GetShiftFile(proj.Source, proj.CloneUrl, proj.Branch)
		if err != nil {
			return errors.Errorf("Failed to get shift file (source: %s, CloneUrl: %s, branch : %s): %v\n", proj.Source, proj.CloneUrl, proj.Branch, err)
		}
	} else {
		f = []byte(proj.GetShiftfile())
	}

	// 5. Parse the shiftfile
	sf, err := parser.AST(f)
	if err != nil {
		return err
	}
	b.f = sf

	m := &types.StorageMetadata{
		TeamID:       b.config.TeamID,
		RepositoryID: proj.GetRepositoryId(),
		BuildID:      b.config.BuildID,
		SubBuildID:   b.config.SubBuildID,
		Branch:       b.project.Branch,
		Path:         b.project.StoragePath,
	}

	// conver storage object
	stor := storage.Convert(proj.Storage)

	// initialize the storage
	ss, err := storage.NewWithMetadata(b.wctx.EnvLogger, stor, m)
	if err != nil {
		return errors.Errorf("Failed to initialize storage: %v", err)
	}
	b.storage = ss

	// start log shipper
	ls, err := logshipper.New(b.wctx.EnvLogger, ss)
	if err != nil {
		return errors.Errorf("Failed to initialize log shipper: %v", err)
	}
	b.logshipper = ls

	b.wctx.EnvTimer.Stop()

	// 6. Ensure the arguments are inputted as static or dynamic values (through env)
	// TODO

	// 7. Construct the runtime execution map from shiftfile ast
	graph, err := graph.Construct(sf)
	if err != nil {
		return err
	}
	graph.SetEnvTimer(b.wctx.EnvTimer)

	b.g = graph

	// 8. Fetch the secrets

	// send the initial graph to server
	b.UpdateBuildGraphToShiftServer("", "", "", b.wctx.EnvLogger)

	// 9. Traverse the execution map & run the actual build
	err = b.build(graph)
	if err != nil {
		b.wctx.EnvLogger.Printf("Build failed: %v\n", err)
	}

	return nil
}
