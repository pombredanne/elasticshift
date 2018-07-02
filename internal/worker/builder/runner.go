/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"log"
	"runtime"
	"sync"
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
		edgeSize := len(c.edges)
		if edgeSize > 0 {

			var errMutex sync.Mutex
			var wg sync.WaitGroup
			parallelCh := make(chan int, parallel)

			for j := 0; j < edgeSize; j++ {

				wg.Add(1)

				go func(n *N) {

					defer wg.Done()
					parallelCh <- 1

					err := b.invokePlugin(n)

					if err != nil {
						errMutex.Lock()
						defer errMutex.Unlock()
						log.Printf("Error when invoking Plugin: %v\n", err)
					}
					<-parallelCh
				}(c.edges[j])
				// s += fmt.Sprintf("(%d) - %s\n", i+1, c.edges[j].Name())
			}

			// wait until all the parallel tasks are finished
			wg.Wait()

		} else {

			// sequential checkpoint execution
			err := b.invokePlugin(c.node)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
