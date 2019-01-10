/*
Copyright 2018 The Elasticshift Authors.
*/
package graph

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/elasticshift/elasticshift/internal/pkg/utils"
	"github.com/elasticshift/shiftfile/ast"
	"github.com/elasticshift/shiftfile/keys"
	"github.com/sirupsen/logrus"
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

	ENV      = "ENV"
	ENV_DESC = "Environment Setup"

	RESTORE_CACHE      = "RCACHE"
	RESTORE_CACHE_DESC = "Restore Cache"

	SAVE_CACHE      = "SCACHE"
	SAVE_CACHE_DESC = "Save Cache"

	ERROR = "ERROR"
)

type FanN struct {
	in  *N
	out *N

	prefix string
	level  int
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
	value map[string]interface{}

	timer utils.Timer

	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	Message     string    `json:"message,omitempty"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	EndedAt     time.Time `json:"ended_at,omitempty"`
	Duration    string    `json:"duration,omitempty"`

	ID string `json:"id,omitempty"`

	Parallel bool
	Logger   *logrus.Entry
}

func newN(value map[string]interface{}) *N {

	n := &N{}
	n.value = value
	n.timer = utils.NewTimer()
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
	if i.Name != ENV {
		i.timer.Start()
		i.StartedAt = i.timer.StartedAt()
	}
}

func (i *N) Wait() {
	if i.Status == "" {
		i.Status = StatusWaiting
	}
}

func (i *N) End(status, message string) {

	if i.Name != ENV {

		i.Status = status
		if message != "" {
			i.Message = message
		}

		i.timer.Stop()
		i.EndedAt = i.timer.StoppedAt()
		i.Duration = i.timer.Duration()
	}
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
		Duration    string    `json:"duration,omitempty"`
		ID          string    `json:"id"`
	}{
		Name:        i.Name,
		Description: i.Description,
		Status:      i.Status,
		StartedAt:   i.StartedAt,
		EndedAt:     i.EndedAt,
		Message:     msg,
		Duration:    i.Duration,
		ID:          i.ID,
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

	prevLevel int
	level     int
	prefix    string

	stack          utils.Stack
	fanoutStack    FanNStack
	fanoutMode     bool
	lastFanoutName string
}

// Construct ...
// Creates a graph with shift file ast.
func Construct(shiftfile *ast.File) (*Graph, error) {

	g := &Graph{
		f:           shiftfile,
		stack:       utils.NewStack(),
		fanoutStack: NewFanNStack(),
		level:       0,
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

func (g *Graph) nodeid() string {

	prefix := g.stack.Last()

	var id string
	if prefix != "" {
		id = prefix + "."
	}

	id += strconv.Itoa(g.level)

	return id
}

func (g *Graph) IncrementPrefix() (string, int) {

	prefix := g.stack.Last()
	if prefix != "" {
		prefix = prefix + "."
	}
	prefix = prefix + strconv.Itoa(g.level-1)
	level := g.level

	g.prevLevel = g.level + 1
	g.level = 1

	g.stack.Push(prefix)

	return prefix, level
}

func (g *Graph) IncrementLevel() {
	g.level++
}

func (g *Graph) DecrementPrefix() {
	g.level = g.prevLevel
	g.prefix = g.stack.Pop()
}

func (g *Graph) constructGraph() error {

	// add start node
	g.addNode(g.constructNode(START, START_DESC))
	g.addNode(g.constructNode(ENV, ENV_DESC))
	g.addNode(g.constructNode(RESTORE_CACHE, RESTORE_CACHE_DESC))

	for g.f.HasMoreBlocks() {

		b := g.f.NextBlock()
		g.addNode(newN(b))
	}

	g.addNode(g.constructNode(SAVE_CACHE, SAVE_CACHE_DESC))

	// add end node
	g.addNode(g.constructNode(END, END_DESC))

	return nil
}

// Update the ENV graph details to existing nodes.
func (g *Graph) SetEnvTimer(timer utils.Timer) {

	n := g.node(ENV)

	n.StartedAt = timer.StartedAt()
	n.EndedAt = timer.StoppedAt()
	n.Duration = timer.Duration()
	n.Status = StatusSuccess
}

func (g *Graph) IsCacheSaved() bool {
	return g.GetSaveCacheNode().Status == StatusSuccess
}

func (g *Graph) GetSaveCacheNode() *N {
	return g.node(SAVE_CACHE)
}

func (g *Graph) node(name string) *N {

	for _, n := range g.nodes {

		if n.Name == name {
			return n
		}
	}
	return nil
}

func (g *Graph) addNode(n *N) {

	g.lock.Lock()

	if name := n.HintName(); name != "" {

		if g.hintOrigins == nil {
			g.hintOrigins = make(map[string]*FanN)
		}

		if g.lastFanoutName != "" && g.lastFanoutName != name {
			g.DecrementPrefix()
		}

		var fann *FanN
		fann = g.hintOrigins[name]
		if fann == nil {

			fann = &FanN{}
			fann.out = g.constructNode(FANOUT+"-"+name, FANOUT_DESC)
			fann.in = g.constructNode(FANIN+"-"+name, FANIN_DESC)

			fann.out.ID = g.nodeid()

			g.IncrementLevel()
			fann.in.ID = g.nodeid()

			// add fan-out node
			g.nodes = append(g.nodes, fann.out)
			g.addEdge(g.prevNode, fann.out)
			g.addCheckpoint(fann.out, n)

			fann.prefix, fann.level = g.IncrementPrefix()

			n.ID = g.nodeid()

			// add fan-in node
			g.nodes = append(g.nodes, n)
			g.addEdge(n, fann.in)
			g.addCheckpoint(fann.in, nil)

			g.hintOrigins[name] = fann
			g.nodes = append(g.nodes, fann.in)
			g.prevNode = fann.in

			g.fanoutStack.Push(fann)
			g.fanoutMode = true

			g.lastFanoutName = name

		} else {

			g.IncrementLevel()
			n.ID = g.nodeid()

			// add edge to a fan-out, fan-in node
			g.nodes = append(g.nodes, n)
			g.addEdge(fann.out, n)
			g.addEdge(n, fann.in)

			g.addCheckpoint(fann.out, n)
		}

	} else {

		if g.fanoutMode {
			g.DecrementPrefix()
			g.fanoutMode = false
		}

		n.ID = g.nodeid()
		g.IncrementLevel()

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

	// s := ""
	// for i := 0; i < len(g.checkpoints); i++ {

	// 	c := g.checkpoints[i]
	// 	s += fmt.Sprintf("(%d) %s\n", i+1, c.Node.Name)
	// 	for j := 0; j < len(c.Edges); j++ {
	// 		s += fmt.Sprintf("(%d) - %s\n", i+1, c.Edges[j].Name)
	// 	}
	// }

	s := ""
	for i := 0; i < len(g.checkpoints); i++ {

		c := g.checkpoints[i]
		s += fmt.Sprintf("(%s) %s\n", c.Node.ID, c.Node.Name)
		for j := 0; j < len(c.Edges); j++ {
			s += fmt.Sprintf("(%s) - %s\n", c.Edges[j].ID, c.Edges[j].Name)
		}
	}

	g.lock.RUnlock()

	return s
}
