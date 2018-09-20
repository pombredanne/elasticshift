/*
Copyright 2018 The Elasticshift Authors.
*/
package graph

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

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
	FANIN_DESC = "Merging the graph to continue further execution"
)

type FanN struct {
	in  *N
	out *N
}

var (
	StatusSuccess    = "S"
	StatusFailed     = "F"
	StatusWaiting    = "W"
	StatusRunning    = "R"
	StatusUnknown    = "U"
	StatusNotStarted = "N"
	StatusCancelled  = "C"
)

// N ...
// Represents a node in execution map
type N struct {
	value map[string]interface{} `json:"-"`

	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	Message     string    `json:"message,omitempty"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	EndedAt     time.Time `json:"ended_at,omitempty"`
}

func newN(value map[string]interface{}) *N {

	n := &N{}
	n.value = value
	n.Name = value[keys.NAME].(string)
	n.Description = value[keys.DESC].(string)
	if n.Name != START || n.Name != END || n.Name != FANOUT || n.Name != FANIN {
		n.Status = StatusNotStarted
	}
	return n
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

func (i *N) SetStatus(status string) {
	i.Status = status
}

func (i *N) String() string {
	return i.Name
}

func (i *N) TimeTaken() string {
	return ""
}

func (i *N) Start() {

	i.Status = StatusRunning
	i.StartedAt = time.Now()
}

func (i *N) Wait() {
	i.Status = StatusWaiting
}

func (i *N) End(status, message string) {
	i.Status = status
	if message != "" {
		i.Message = message
	}
	i.EndedAt = time.Now()
}

// MarshalJSON ..
// Serialize the N (node) to json format
func (i *N) MarshalJSON() ([]byte, error) {

	var msg string
	if i.Message != "" {
		msg = base64.StdEncoding.EncodeToString([]byte(i.Message))
	}

	return json.Marshal(&struct {
		Name        string    `json:"name"`
		Description string    `json:"description,omitempty"`
		Status      string    `json:"status,omitempty"`
		Message     string    `json:"message,omitempty"`
		StartedAt   time.Time `json:"started_at,omitempty"`
		EndedAt     time.Time `json:"ended_at,omitempty"`
	}{
		Name:        i.Name,
		Description: i.Description,
		Status:      i.Status,
		StartedAt:   i.StartedAt,
		EndedAt:     i.EndedAt,
		Message:     msg,
	})
}

// Checkpoint ..
// Each and every hop of the build during execution.
type Checkpoint struct {
	Node  *N   `json:"node"`
	Edges []*N `json:"edges,omitempty"`
}

// Graph ..
// Shiftfile is representated in graph format
type Graph struct {
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

// Construct ...
// Creates a graph with shift file ast.
func Construct(shiftfile *ast.File) (*Graph, error) {

	g := &Graph{
		f: shiftfile,
	}

	err := g.constructGraph()
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Graph) constructNode(name, description string) *N {

	v := make(map[string]interface{})
	v[keys.NAME] = name
	v[keys.DESC] = description

	return newN(v)
}

func (g *Graph) constructGraph() error {

	// add start node
	g.addNode(g.constructNode(START, START_DESC))

	for g.f.HasMoreBlocks() {

		b := g.f.NextBlock()
		g.addNode(newN(b))
	}

	// add end node
	g.addNode(g.constructNode(END, END_DESC))

	return nil
}

func (g *Graph) addNode(n *N) {

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

func (g *Graph) addEdge(n1, n2 *N) {

	if g.edges == nil {
		g.edges = make(map[*N][]*N)
	}
	g.edges[n1] = append(g.edges[n1], n2)
}

func (g *Graph) addCheckpoint(n *N, e *N) {

	if n == e {
		return
	}

	if g.checkpoints == nil {
		g.checkpoints = make([]*Checkpoint, 0)
	}

	var cp *Checkpoint
	var found bool
	for i := 0; i < len(g.checkpoints); i++ {

		if n == g.checkpoints[i].Node {

			cp = g.checkpoints[i]
			found = true
			break
		}
	}

	if !found {

		cp = &Checkpoint{}

		cp.Edges = make([]*N, 0)

		g.checkpoints = append(g.checkpoints, cp)
	}

	cp.Node = n

	if e != nil {
		cp.Edges = append(cp.Edges, e)
	}
}

// Checkpoints ...
// Gets the checpoints
func (g *Graph) Checkpoints() []*Checkpoint {
	return g.checkpoints
}

// JSON ...
// Return graph in json format
func (g *Graph) JSON() (string, error) {

	g.lock.RLock()

	nods, err := json.Marshal(g.checkpoints)
	if err != nil {
		return "", err
	}

	g.lock.RUnlock()

	return string(nods), nil
}

func (g *Graph) String() string {

	g.lock.RLock()

	s := ""
	for i := 0; i < len(g.checkpoints); i++ {

		c := g.checkpoints[i]
		s += fmt.Sprintf("(%d) %s\n", i+1, c.Node.Name)
		for j := 0; j < len(c.Edges); j++ {
			s += fmt.Sprintf("(%d) - %s\n", i+1, c.Edges[j].Name)
		}
	}
	g.lock.RUnlock()

	return s
}