/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"log"
	"os"
	"runtime"
	"sync"

	homedir "github.com/minio/go-homedir"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/internal/pkg/graph"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
)

func (b *builder) build(g *graph.Graph) error {

	wdir := b.f.WorkDir()

	if wdir != "" {
		expanded, err := homedir.Expand(wdir)
		if err != nil {
			log.Printf("Failed to expand the directory : %v\n", err)
		}

		err = utils.Mkdir(expanded)
		if err != nil {
			log.Printf("Failed to create working directory : %v\n", err)
		}

		err = os.Chdir(expanded)
		if err != nil {
			log.Printf("Failed to change the working directory : %v\n", err)
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
			log.Println("S:FO:~" + c.Node.ID + "~")

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

					b.UpdateBuildGraphToShiftServer(graph.StatusWaiting, n.Name, "")

					defer wg.Done()
					parallelCh <- 1

					n.Start()
					b.UpdateBuildGraphToShiftServer(graph.StatusRunning, n.Name, "")

					msg, err := b.invokePlugin(n)

					if err != nil {

						errMutex.Lock()
						defer errMutex.Unlock()

						log.Printf("Failed when executing %s : %v\n", n.Name, err)
						n.End(graph.StatusFailed, msg)
						b.logBlockInfo(n, "F")
						b.UpdateBuildGraphToShiftServer(graph.StatusFailed, n.Name, msg)

						failed = true
					} else {

						if n.Status != graph.StatusFailed {
							n.End(graph.StatusSuccess, "")
							b.logBlockInfo(n, "E")
							b.UpdateBuildGraphToShiftServer(graph.StatusSuccess, n.Name, "")
						}
					}

					<-parallelCh

				}(c.Edges[j])
			}

			// wait until all the parallel tasks are finished
			wg.Wait()

			if fanoutFailed {
				c.Node.End(graph.StatusFailed, "")
				b.logBlockInfo(c.Node, "F")
			} else {
				c.Node.End(graph.StatusSuccess, "")
				b.logBlockInfo(c.Node, "E")
			}
			log.Println("E:FO:~" + c.Node.ID + ":" + c.Node.Duration + "~")

		} else {

			c.Node.Start()
			b.UpdateBuildGraphToShiftServer(graph.StatusRunning, c.Node.Name, "")

			// sequential checkpoint execution
			msg, err := b.invokePlugin(c.Node)
			if err != nil {
				c.Node.End(graph.StatusFailed, msg)

				b.logBlockInfo(c.Node, "F")
				b.UpdateBuildGraphToShiftServer(graph.StatusFailed, c.Node.Name, msg)

				failed = true
			} else {

				c.Node.End(graph.StatusSuccess, "")

				if graph.START != c.Node.Name || graph.END != c.Node.Name {
					b.logBlockInfo(c.Node, "E")
				}
				b.UpdateBuildGraphToShiftServer(graph.StatusSuccess, c.Node.Name, "")
			}
		}
	}

	return nil
}

func (b *builder) UpdateBuildGraphToShiftServer(status, checkpoint, reason string) {

	if graph.StatusFailed == status || (graph.END == checkpoint && graph.StatusSuccess == status) {

		log.Println("S:~0.3:Saving cache:~")

		b.saveCache()

		// TODO add duration
		log.Println("E:~0.3:Saving cache::~")
	}

	gph, err := b.g.JSON()
	if err != nil {
		log.Printf("Eror when contructing status graph: %v", err)
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
			log.Printf("Failed to update buld graph: %v", err)
		}
	}
}
