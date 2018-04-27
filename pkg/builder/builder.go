/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"context"
	"fmt"
	"log"

	"gitlab.com/conspico/elasticshift/api"
	wtypes "gitlab.com/conspico/elasticshift/pkg/worker/types"
	"google.golang.org/grpc"
)

type builder struct {
	shiftconn   *grpc.ClientConn
	ctx         context.Context
	config      wtypes.Config
	shiftclient api.ShiftClient
	project     *api.GetProjectRes
}

func New(ctx wtypes.Context, shiftconn *grpc.ClientConn) error {

	b := builder{}
	b.shiftconn = shiftconn
	b.ctx = ctx.Context
	b.shiftclient = ctx.Client
	b.config = ctx.Config

	return b.run()
}

func (b *builder) run() error {

	// Get the project information
	proj, err := b.shiftclient.GetProject(b.ctx, &api.GetProjectReq{BuildId: b.config.BuildID})
	if err != nil {
		return fmt.Errorf("Failed to get the project/repository detail from shift server: %v", err)
	}
	b.project = proj

	log.Printf("Project Info: %v", proj)

	// Checkout the source code

	// Analyze the build spec (shiftfile), if exist
	// otherwise use the global language spec defined by elasticshift

	return nil
}
