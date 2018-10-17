/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/Sirupsen/logrus"
	homedir "github.com/minio/go-homedir"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/internal/pkg/graph"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
)

func (b *builder) build(g *graph.Graph) error {

	log1 := b.wctx.EnvLogger

	wdir := b.f.WorkDir()

	if wdir != "" {
		expanded, err := homedir.Expand(wdir)
		if err != nil {
			log1.Printf("Failed to expand the directory : %v\n", err)
		}

		err = utils.Mkdir(expanded)
		if err != nil {
			log1.Printf("Failed to create working directory : %v\n", err)
		}

		err = os.Chdir(expanded)
		if err != nil {
			log1.Printf("Failed to change the working directory : %v\n", err)
		}
	}

	// set the parallel capability
	var parallel int
	nCpu := runtime.NumCPU()
	if nCpu < 2 {
		parallel = 1
	} else {
		parallel = nCpu - 1
	}

	var failed bool

	// walk through the checkpoints and execute them
	// s := ""
	checkpoints := g.Checkpoints()
	for i := 0; i < len(checkpoints); i++ {

		if failed {
			b.done <- 1
			return nil
		}

		c := checkpoints[i]
		// s += fmt.Sprintf("(%d) %s\n", i+1, c.node.Name())

		// run block if it is a sequential task
		// If it's FANOUT, spawn the multiple block
		// with in a worker group and wait for it to complete
		edgeSize := len(c.Edges)
		if edgeSize > 0 {

			var fanoutFailed bool

			c.Node.Wait()

			var errMutex sync.Mutex
			var wg sync.WaitGroup
			parallelCh := make(chan int, parallel)

			for j := 0; j < edgeSize; j++ {

				if failed {
					fanoutFailed = true
					break
				}

				wg.Add(1)

				if j == 0 {
					c.Node.Start()
				}

				go func(n *graph.N) {

					n.Wait()
					n.Parallel = true

					nodelogger, err := b.wctx.LogWriter.GetLogger(n.ID)
					if err != nil {
						fmt.Printf("Error when getting logger for nodeid : %s", n.ID)
					}
					n.Logger = nodelogger

					b.UpdateBuildGraphToShiftServer(graph.StatusWaiting, n.Name, "", nodelogger)

					defer wg.Done()
					parallelCh <- 1

					n.Start()
					b.UpdateBuildGraphToShiftServer(graph.StatusRunning, n.Name, "", nodelogger)

					msg, err := b.invokePlugin(n)

					if err != nil {

						errMutex.Lock()
						defer errMutex.Unlock()

						nodelogger.Printf("Failed when executing %s : %v\n", n.Name, err)
						n.End(graph.StatusFailed, msg)
						b.UpdateBuildGraphToShiftServer(graph.StatusFailed, n.Name, msg, nodelogger)

						failed = true
					} else {

						if n.Status != graph.StatusFailed {
							n.End(graph.StatusSuccess, "")
							b.UpdateBuildGraphToShiftServer(graph.StatusSuccess, n.Name, "", nodelogger)
						}
					}

					<-parallelCh

				}(c.Edges[j])
			}

			// wait until all the parallel tasks are finished
			wg.Wait()

			if fanoutFailed {
				c.Node.End(graph.StatusFailed, "")
			} else {
				c.Node.End(graph.StatusSuccess, "")
			}

		} else {

			nodelogger, err := b.wctx.LogWriter.GetLogger(c.Node.ID)
			if err != nil {
				fmt.Printf("Error when getting logger for nodeid : %s", c.Node.ID)
			}
			c.Node.Logger = nodelogger

			c.Node.Start()
			b.UpdateBuildGraphToShiftServer(graph.StatusRunning, c.Node.Name, "", nodelogger)

			// sequential checkpoint execution
			msg, err := b.invokePlugin(c.Node)
			if err != nil {
				c.Node.End(graph.StatusFailed, msg)
				b.UpdateBuildGraphToShiftServer(graph.StatusFailed, c.Node.Name, msg, nodelogger)

				failed = true
			} else {

				c.Node.End(graph.StatusSuccess, "")
				b.UpdateBuildGraphToShiftServer(graph.StatusSuccess, c.Node.Name, "", nodelogger)
			}
		}
	}

	return nil
}

func (b *builder) UpdateBuildGraphToShiftServer(status, checkpoint, reason string, logn *logrus.Entry) {

	// if graph.StatusFailed == status || (graph.END == checkpoint && graph.StatusSuccess == status) {

	// 	log.Println("S:~0.3:Saving cache:~")

	// 	b.saveCache()

	// 	// TODO add duration
	// 	log.Println("E:~0.3:Saving cache::~")
	// }

	gph, err := b.g.JSON()
	if err != nil {
		logn.Printf("Eror when contructing status graph: %v\n", err)
	}

	req := &api.UpdateBuildStatusReq{}
	req.BuildId = b.config.BuildID
	req.Graph = gph
	req.Status = status
	req.Checkpoint = checkpoint
	if reason != "" {
		req.Reason = reason
	}

	if b.shiftclient != nil {
		_, err = b.shiftclient.UpdateBuildStatus(b.ctx, req)
		if err != nil {
			logn.Printf("Failed to update buld graph: %v\n", err)
		}
	}
}
