/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"fmt"
	"sync"
	"time"

	"encoding/json"

	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/ast"
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/keys"
)

var (
	START      = "START"
	START_DESC = "Starting the execution/graph"

	END      = "END"
	END_DESC = "Finishing the execution/graph"

	FANOUT      = "FANOUT"
	FANOUT_DESC = "Spitting the graph for parallel execution"

	FANIN      = "FANIN"
	FANIN_DESC = "Merging the graph to continue sequential execution"
)

type FanN struct {
	in  *N
	out *N
}

type N struct {
	value map[string]interface{}

	start time.Time
	end   time.Time
}

func (i *N) Name() string {
	return i.value[keys.NAME].(string)
}

func (i *N) Description() string {
	return i.value[keys.DESC].(string)
}

func (i *N) Item() map[string]interface{} {
	return i.value
}

func (i *N) HintName() string {
	hmap := i.value[keys.HINT]
	if hmap != nil {
		return hmap.(map[string]string)["PARALLEL"]
	}
	return ""
}

func (i *N) String() string {
	return i.Name()
}

func (i *N) StartedAt() string {
	return i.start.String()
}

func (i *N) EndedAt() string {
	return i.end.String()
}

func (i *N) TimeTaken() string {
	return ""
}

type Checkpoint struct {
	node  *N
	edges []*N
}

type graph struct {
	f *ast.File

	nodes       []*N
	checkpoints []*Checkpoint

	edges map[*N][]*N

	startNode *N
	endNode   *N

	prevNode *N

	hintOrigins map[string]*FanN

	lock sync.RWMutex
}

func ConstructGraph(shiftfile *ast.File) (*graph, error) {

	g := &graph{
		f: shiftfile,
	}

	err := g.constructGraph()
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *graph) constructNode(name, description string) *N {

	v := make(map[string]interface{})
	v[keys.NAME] = name
	v[keys.DESC] = description

	return &N{value: v}
}

func (g *graph) constructGraph() error {

	// add start node
	g.addNode(g.constructNode(START, START_DESC))

	for g.f.HasMoreBlocks() {

		b := g.f.NextBlock()
		g.addNode(&N{value: b})
	}

	// add end node
	g.addNode(g.constructNode(END, END_DESC))

	return nil
}

func (g *graph) addNode(n *N) {

	g.lock.Lock()

	if name := n.HintName(); name != "" {

		if g.hintOrigins == nil {
			g.hintOrigins = make(map[string]*FanN)
		}

		var fann *FanN
		fann = g.hintOrigins[name]
		if fann == nil {

			fann = &FanN{}
			fann.out = g.constructNode(FANOUT+"-"+name, FANOUT_DESC)
			fann.in = g.constructNode(FANIN+"-"+name, FANIN_DESC)

			// add fan-out node
			g.nodes = append(g.nodes, fann.out)
			g.addEdge(g.prevNode, fann.out)
			g.addCheckpoint(fann.out, n)

			// add fan-in node
			g.nodes = append(g.nodes, n)
			g.addEdge(n, fann.in)
			g.addCheckpoint(fann.in, nil)

			g.hintOrigins[name] = fann
			g.nodes = append(g.nodes, fann.in)
			g.prevNode = fann.in

		} else {

			// add edge to a fan-out, fan-in node
			g.nodes = append(g.nodes, n)
			g.addEdge(fann.out, n)
			g.addEdge(n, fann.in)

			g.addCheckpoint(fann.out, n)
		}

	} else {

		g.nodes = append(g.nodes, n)
		g.prevNode = n

		g.addCheckpoint(n, nil)
	}

	g.lock.Unlock()
}

func (g *graph) addEdge(n1, n2 *N) {

	if g.edges == nil {
		g.edges = make(map[*N][]*N)
	}
	g.edges[n1] = append(g.edges[n1], n2)
}

func (g *graph) addCheckpoint(n *N, e *N) {

	if n == e {
		return
	}

	if g.checkpoints == nil {
		g.checkpoints = make([]*Checkpoint, 0)
	}

	var cp *Checkpoint
	var found bool
	for i := 0; i < len(g.checkpoints); i++ {

		if n == g.checkpoints[i].node {

			cp = g.checkpoints[i]
			found = true
			break
		}
	}

	if !found {

		cp = &Checkpoint{}

		cp.edges = make([]*N, 0)

		g.checkpoints = append(g.checkpoints, cp)
	}

	cp.node = n

	if e != nil {
		cp.edges = append(cp.edges, e)
	}
}

func (g *graph) Checkpoints() []*Checkpoint {
	return g.checkpoints
}

func (g *graph) Json() string {

	nods, _ := json.Marshal(g.nodes)
	return string(nods)
}

func (g *graph) String() string {

	g.lock.RLock()

	s := ""
	for i := 0; i < len(g.checkpoints); i++ {

		c := g.checkpoints[i]
		s += fmt.Sprintf("(%d) %s\n", i+1, c.node.Name())
		for j := 0; j < len(c.edges); j++ {
			s += fmt.Sprintf("(%d) - %s\n", i+1, c.edges[j].Name())
		}
	}
	g.lock.RUnlock()

	return s
}
