/*
Copyright 2018 The Elasticshift Authors.
*/
package logshipper

import (
	"fmt"
	"sync"

	"github.com/golang/protobuf/ptypes"
	"gitlab.com/conspico/elasticshift/api"
	stypes "gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/pkg/worker/types"
)

type embedded struct {
	mutx   sync.RWMutex
	ctx    types.Context
	logs   chan stypes.Log
	stream api.Shift_LogShipClient
}

func (l *embedded) Send(log stypes.Log) {
	l.logs <- log
}

func (l *embedded) Log(msg string) {
	log := constructLog(msg)
	l.Send(log)
}

func (l *embedded) Error(err error) {
	log := constructLog(err.Error())
	l.Send(log)
}

func (l *embedded) Info(msg string) {
	log := constructLog(msg)
	l.Send(log)
}

func (l *embedded) Halt() error {
	return l.stream.CloseSend()
}

func newEmbeddedLogger(ctx types.Context) (Logger, error) {

	logr := &embedded{
		ctx:  ctx,
		logs: make(chan stypes.Log),
	}

	stream, err := ctx.Client.LogShip(ctx.Context)
	if err != nil {
		return nil, err
	}

	logr.stream = stream

	go func(stream api.Shift_LogShipClient, ctx types.Context) {

		for log := range logr.logs {

			req := &api.LogShipReq{
				BuildId: ctx.Config.BuildID,
				Log:     log.Data,
			}
			req.Time, _ = ptypes.TimestampProto(log.Time)

			logr.mutx.RLock()

			// Send the log to elasticshift server
			err := stream.Send(req)
			if err != nil {
				fmt.Println("Can't send message: ", err.Error())
			}
			logr.mutx.RUnlock()

		}
	}(stream, ctx)

	logr.stream = stream
	return logr, nil
}
