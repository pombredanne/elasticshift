/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"context"
	"fmt"
	"io"
	"log"

	"path/filepath"

	"github.com/pkg/errors"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/internal/pkg/graph"
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/ast"
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/parser"
	"gitlab.com/conspico/elasticshift/internal/pkg/vcs"
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
	config      wtypes.Config
	shiftclient api.ShiftClient
	project     *api.GetProjectRes

	f *ast.File
	g *graph.Graph

	done chan int

	writer io.Writer
}

func New(ctx wtypes.Context, shiftconn *grpc.ClientConn, writer io.Writer, done chan int) error {

	b := builder{}
	b.shiftconn = shiftconn
	b.ctx = ctx.Context
	b.shiftclient = ctx.Client
	b.config = ctx.Config
	b.done = done
	b.writer = writer

	return b.run()
}

func (b *builder) run() error {

	// restore build cache if any
	// save the cache after every successful build
	err := b.restoreCache()
	if err != nil {
		log.Println("Restoring cache failed:", err)
	}

	// Get the project information
	proj, err := b.shiftclient.GetProject(b.ctx, &api.GetProjectReq{BuildId: b.config.BuildID})
	if err != nil {
		return fmt.Errorf("Failed to get the project/repository detail from shift server: %v", err)
	}
	b.project = proj

	log.Printf("Project Info: %v", proj)

	// 1. Ensure connection to log storage is good, this container should be loaded with

	// 2. Load the build cache, if available ensure it

	// 3. Fetch the shiftfile

	log.Println("Getting the shift file..")

	f, err := vcs.GetShiftFile(proj.Source, proj.CloneUrl, proj.Branch)
	if err != nil {
		return errors.Errorf("Failed to get shift file (source: %s, CloneUrl: %s, branch : %s): %v", proj.Source, proj.CloneUrl, proj.Branch, err)
	}

	// 4. otherwise use the global language spec defined by elasticshift
	if f == nil {
		//TODO fetch the default shift file
	}

	// 5. Parse the shiftfile
	sf, err := parser.AST(f)
	if err != nil {
		return err
	}
	b.f = sf

	// 6. Ensure the arguments are inputted as static or dynamic values (through env)

	// 7. Construct the runtime execution map from shiftfile ast
	graph, err := graph.Construct(sf)
	if err != nil {
		return err
	}
	b.g = graph

	// 8. Fetch the secrets

	// send the initial graph to server
	b.UpdateBuildGraphToShiftServer("", "")

	// 9. Traverse the execution map & run the actual build
	err = b.build(graph)
	if err != nil {
		log.Printf("Build failed: %v\n", err)
	}

	return nil
}
