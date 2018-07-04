/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"log"
	"runtime"
	"sync"

	"gitlab.com/conspico/elasticshift/api"
)

func (b *builder) build(g *graph) error {

	// set the parallel capability
	var parallel int
	nCpu := runtime.NumCPU()
	if nCpu < 2 {
		parallel = 1
	} else {
		parallel = nCpu - 1
	}

	// walk through the checkpoints and execute them
	// s := ""
	for i := 0; i < len(g.checkpoints); i++ {

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

				wg.Add(1)

				go func(n *N) {

					n.Wait()
					b.UpdateBuildGraphToShiftServer(statusWaiting, n.Name)

					defer wg.Done()
					parallelCh <- 1

					n.Start()
					b.UpdateBuildGraphToShiftServer(statusRunning, n.Name)

					err := b.invokePlugin(n)

					if err != nil {
						errMutex.Lock()
						defer errMutex.Unlock()
						log.Printf("Error when invoking Plugin: %v\n", err)
						n.End(statusFailed, err.Error())
						b.UpdateBuildGraphToShiftServer(statusFailed, n.Name)
					}

					if n.Status != statusFailed {
						n.End(statusSuccess, "")
						b.UpdateBuildGraphToShiftServer(statusSuccess, n.Name)
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
			err := b.invokePlugin(c.Node)
			if err != nil {
				c.Node.End(statusFailed, err.Error())
				b.UpdateBuildGraphToShiftServer(statusFailed, c.Node.Name)
				return err
			}

			c.Node.End(statusSuccess, "")
			b.UpdateBuildGraphToShiftServer(statusSuccess, c.Node.Name)
		}
	}

	// finishes the build
	b.done <- 1

	return nil
}

func (b *builder) UpdateBuildGraphToShiftServer(status, checkpoint string) {

	gph, err := b.g.Json()
	if err != nil {
		log.Println("Eror when contructing status graph: %v", err)
	}

	req := &api.UpdateBuildStatusReq{}
	req.BuildId = b.config.BuildID
	req.Graph = gph
	req.Status = status
	req.Checkpoint = checkpoint

	_, err = b.shiftclient.UpdateBuildStatus(b.ctx, req)
	if err != nil {
		log.Println("Failed to update buld graph: %v", err)
	}
}
