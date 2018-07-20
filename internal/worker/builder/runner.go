/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"os"
	"runtime"
	"sync"

	homedir "github.com/minio/go-homedir"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
)

func (b *builder) build(g *graph) error {

	wdir := b.f.WorkDir()
	b.log.Debugf("Working directory : %s\n", wdir)

	if wdir != "" {
		expanded, err := homedir.Expand(wdir)
		if err != nil {
			b.log.Errorf("Failed to expand the directory : %v\n", err)
		}

		err = utils.Mkdir(expanded)
		if err != nil {
			b.log.Errorf("Failed to create working directory : %v\n", err)
		}

		err = os.Chdir(expanded)
		if err != nil {
			b.log.Errorf("Failed to change the working directory : %v\n", err)
		}
	}

	b.logr.Log("Working directory = " + utils.GetWD())

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
	for i := 0; i < len(g.checkpoints); i++ {

		if failed {
			b.log.Infoln("Build finished, waiting from shiftserver to receive a halt command..")
			<-b.done
		}

		c := g.checkpoints[i]
		// s += fmt.Sprintf("(%d) %s\n", i+1, c.node.Name())

		// run block if it is a sequential task
		// If it's FANOUT, spawn the multiple block
		// with in a worker group and wait for it to complete
		edgeSize := len(c.Edges)
		if edgeSize > 0 {

			var errMutex sync.Mutex
			var wg sync.WaitGroup
			parallelCh := make(chan int, parallel)

			for j := 0; j < edgeSize; j++ {

				if failed {
					break
				}

				wg.Add(1)

				go func(n *N) {

					n.Wait()
					b.UpdateBuildGraphToShiftServer(statusWaiting, n.Name)

					defer wg.Done()
					parallelCh <- 1

					n.Start()
					b.UpdateBuildGraphToShiftServer(statusRunning, n.Name)

					msg, err := b.invokePlugin(n)

					if err != nil {

						errMutex.Lock()
						defer errMutex.Unlock()

						b.log.Errorf("Plugin error : %v\n", err)
						n.End(statusFailed, msg)
						b.UpdateBuildGraphToShiftServer(statusFailed, n.Name)

						failed = true
					} else {

						if n.Status != statusFailed {
							n.End(statusSuccess, "")
							b.UpdateBuildGraphToShiftServer(statusSuccess, n.Name)
						}
					}

					<-parallelCh

				}(c.Edges[j])
			}

			// wait until all the parallel tasks are finished
			wg.Wait()

		} else {

			c.Node.Start()
			b.UpdateBuildGraphToShiftServer(statusRunning, c.Node.Name)

			// sequential checkpoint execution
			msg, err := b.invokePlugin(c.Node)
			if err != nil {
				c.Node.End(statusFailed, msg)
				b.log.Errorf("Plugin error : %v\n", err)
				b.UpdateBuildGraphToShiftServer(statusFailed, c.Node.Name)

				failed = true
			} else {

				c.Node.End(statusSuccess, "")
				b.UpdateBuildGraphToShiftServer(statusSuccess, c.Node.Name)
			}
		}
	}

	return nil
}

func (b *builder) UpdateBuildGraphToShiftServer(status, checkpoint string) {

	if statusFailed == status || (END == checkpoint && statusSuccess == status) {

		b.log.Infoln("Saving cache.")

		b.saveCache()

		b.log.Infoln("Finished saving the cache")
	}

	gph, err := b.g.Json()
	if err != nil {
		b.log.Errorf("Eror when contructing status graph: %v", err)
	}

	req := &api.UpdateBuildStatusReq{}
	req.BuildId = b.config.BuildID
	req.Graph = gph
	req.Status = status
	req.Checkpoint = checkpoint

	if b.shiftclient != nil {
		_, err = b.shiftclient.UpdateBuildStatus(b.ctx, req)
		if err != nil {
			b.log.Errorf("Failed to update buld graph: %v", err)
		}
	}
}
