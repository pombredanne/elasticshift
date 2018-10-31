/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	homedir "github.com/minio/go-homedir"
	"github.com/sirupsen/logrus"
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
	parallel := utils.NumOfCPU()

	var failed bool

	// walk through the checkpoints and execute them
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

					defer b.ShipLog(n.ID, n.Name)
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
			failed = b.runNode(c.Node)
		}
	}

	return nil
}

func (b *builder) runNode(n *graph.N) bool {

	var failed bool

	nodelogger, err := b.wctx.LogWriter.GetLogger(n.ID)
	if err != nil {
		fmt.Printf("Error when getting logger for nodeid : %s", n.ID)
	}
	n.Logger = nodelogger

	// defer b.ShipLog(n.ID, n.Name)

	n.Start()
	b.UpdateBuildGraphToShiftServer(graph.StatusRunning, n.Name, "", nodelogger)

	// sequential checkpoint execution
	msg, err := b.invokePlugin(n)
	if err != nil {
		n.End(graph.StatusFailed, msg)

		b.ShipLog(n.ID, n.Name)
		b.UpdateBuildGraphToShiftServer(graph.StatusFailed, n.Name, msg, nodelogger)

		failed = true
	} else {

		n.End(graph.StatusSuccess, "")

		b.ShipLog(n.ID, n.Name)
		b.UpdateBuildGraphToShiftServer(graph.StatusSuccess, n.Name, "", nodelogger)
	}

	return failed
}

func (b *builder) ShipLog(nodeid, name string) {

	if name == graph.START || name == graph.END || strings.HasPrefix(name, graph.FANOUT) || strings.HasPrefix(name, graph.FANIN) {
		return
	}

	lpath, _ := b.wctx.LogWriter.LogPath(nodeid)
	f, _ := b.wctx.LogWriter.LogFile(nodeid)
	f.Close()

	b.logshipper.Ship(nodeid, lpath)
}

func (b *builder) UpdateBuildGraphToShiftServer(status, checkpoint, reason string, logn *logrus.Entry) {

	req := &api.UpdateBuildStatusReq{}
	if graph.StatusFailed == status || (graph.END == checkpoint && graph.StatusSuccess == status) {

		if !b.g.IsCacheSaved() {
			n := b.g.GetSaveCacheNode()
			b.runNode(n)
		}

		req.Duration = utils.CalculateDuration(b.wctx.EnvTimer.StartedAt(), time.Now())

		// wait until log shipper finish
		b.logshipper.WaitUntilLogShipperCompletes()
	}

	gph, err := b.g.JSON()
	if err != nil {
		logn.Printf("Eror when contructing status graph: %v\n", err)
	}

	//fmt.Printf("Config = %#v", b.wctx.Config)
	req.BuildId = b.wctx.Config.BuildID
	req.SubBuildId = b.wctx.Config.SubBuildID
	req.TeamId = b.wctx.Config.TeamID
	req.RepositoryId = b.project.GetRepositoryId()
	req.Branch = b.project.GetBranch()
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
