package orchestrate

import (
	"fmt"
	"sync"

	"github.com/lingzhi227/go-scheduler/pkg/agent"
)

// NodeType defines the type of a graph node.
type NodeType string

const (
	NodeTypeStart    NodeType = "start"
	NodeTypeAgent    NodeType = "agent"
	NodeTypeTool     NodeType = "tool"
	NodeTypeDecision NodeType = "decision"
	NodeTypeJoin     NodeType = "join"
	NodeTypeEnd      NodeType = "end"
)

// NodeStatus tracks execution state.
type NodeStatus string

const (
	NodePending  NodeStatus = "pending"
	NodeRunning  NodeStatus = "running"
	NodeDone     NodeStatus = "done"
	NodeFailed   NodeStatus = "failed"
	NodeSkipped  NodeStatus = "skipped"
)

// Node represents a single step in a TaskGraph.
type Node struct {
	ID          string             `json:"id"`
	Type        NodeType           `json:"type"`
	Label       string             `json:"label"`
	AgentConfig *agent.AgentConfig `json:"agent_config,omitempty"`
	ToolName    string             `json:"tool_name,omitempty"`
	Input       map[string]any     `json:"input,omitempty"`
}

// Edge connects two nodes, optionally with a condition.
type Edge struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Condition string `json:"condition,omitempty"` // optional condition expression
}

// TaskGraph represents a directed acyclic graph of tasks.
type TaskGraph struct {
	mu    sync.RWMutex
	ID    string            `json:"id"`
	Goal  string            `json:"goal"`
	Nodes map[string]*Node  `json:"nodes"`
	Edges []*Edge           `json:"edges"`
	State map[string]any    `json:"state,omitempty"` // shared execution state
}

func NewTaskGraph(id, goal string) *TaskGraph {
	return &TaskGraph{
		ID:    id,
		Goal:  goal,
		Nodes: make(map[string]*Node),
		State: make(map[string]any),
	}
}

// AddNode adds a node to the graph.
func (g *TaskGraph) AddNode(n *Node) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.Nodes[n.ID]; exists {
		return fmt.Errorf("node %q already exists", n.ID)
	}
	g.Nodes[n.ID] = n
	return nil
}

// AddEdge adds a directed edge between two nodes.
func (g *TaskGraph) AddEdge(from, to string) error {
	return g.AddConditionalEdge(from, to, "")
}

// AddConditionalEdge adds a directed edge with an optional condition.
func (g *TaskGraph) AddConditionalEdge(from, to, condition string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.Nodes[from]; !ok {
		return fmt.Errorf("source node %q not found", from)
	}
	if _, ok := g.Nodes[to]; !ok {
		return fmt.Errorf("target node %q not found", to)
	}
	g.Edges = append(g.Edges, &Edge{From: from, To: to, Condition: condition})
	return nil
}

// Predecessors returns the IDs of all predecessors of a node.
func (g *TaskGraph) Predecessors(nodeID string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var preds []string
	for _, e := range g.Edges {
		if e.To == nodeID {
			preds = append(preds, e.From)
		}
	}
	return preds
}

// Successors returns the IDs of all successors of a node.
func (g *TaskGraph) Successors(nodeID string) []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var succs []string
	for _, e := range g.Edges {
		if e.From == nodeID {
			succs = append(succs, e.To)
		}
	}
	return succs
}

// TopologicalSort returns nodes in topological order.
// Returns an error if the graph contains a cycle.
func (g *TaskGraph) TopologicalSort() ([][]string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Compute in-degrees
	inDegree := make(map[string]int)
	for id := range g.Nodes {
		inDegree[id] = 0
	}
	for _, e := range g.Edges {
		inDegree[e.To]++
	}

	// Find nodes with zero in-degree
	var queue []string
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	var levels [][]string
	visited := 0

	for len(queue) > 0 {
		levels = append(levels, queue)
		visited += len(queue)

		var next []string
		for _, id := range queue {
			for _, e := range g.Edges {
				if e.From == id {
					inDegree[e.To]--
					if inDegree[e.To] == 0 {
						next = append(next, e.To)
					}
				}
			}
		}
		queue = next
	}

	if visited != len(g.Nodes) {
		return nil, fmt.Errorf("graph contains a cycle (visited %d of %d nodes)", visited, len(g.Nodes))
	}
	return levels, nil
}

// SetState sets a key in the shared execution state.
func (g *TaskGraph) SetState(key string, value any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.State[key] = value
}

// GetState reads a key from the shared execution state.
func (g *TaskGraph) GetState(key string) (any, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	v, ok := g.State[key]
	return v, ok
}
