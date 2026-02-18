package comm

import (
	"math/rand"
	"sort"
)

// Topology defines which agents can communicate with each other.
// Edges are directed: edges[A] = [B, C] means A can send to B and C.
type Topology struct {
	agents []string
	edges  map[string][]string
}

// NewFullTopology creates a fully connected topology (all-to-all).
func NewFullTopology(agents []string) *Topology {
	t := &Topology{agents: agents, edges: make(map[string][]string)}
	for _, a := range agents {
		for _, b := range agents {
			if a != b {
				t.edges[a] = append(t.edges[a], b)
			}
		}
	}
	return t
}

// NewRingTopology creates a ring where each agent can send to its successor.
func NewRingTopology(agents []string) *Topology {
	t := &Topology{agents: agents, edges: make(map[string][]string)}
	n := len(agents)
	for i, a := range agents {
		next := agents[(i+1)%n]
		t.edges[a] = []string{next}
	}
	return t
}

// NewStarTopology creates a star with a central hub and spokes.
// The hub can communicate with all spokes and vice versa.
func NewStarTopology(hub string, spokes []string) *Topology {
	t := &Topology{
		agents: append([]string{hub}, spokes...),
		edges:  make(map[string][]string),
	}
	t.edges[hub] = make([]string, len(spokes))
	copy(t.edges[hub], spokes)
	for _, s := range spokes {
		t.edges[s] = []string{hub}
	}
	return t
}

// NewSparseTopology creates a random sparse topology with the given connectivity ratio.
// connectivity=0.3 means each agent connects to ~30% of others.
// Paper reference: Sparse Communication Topology (EMNLP 2024).
func NewSparseTopology(agents []string, connectivity float64) *Topology {
	t := &Topology{agents: agents, edges: make(map[string][]string)}
	n := len(agents)
	if n <= 1 {
		return t
	}

	for _, a := range agents {
		others := make([]string, 0, n-1)
		for _, b := range agents {
			if a != b {
				others = append(others, b)
			}
		}
		// Shuffle and pick connectivity fraction
		rand.Shuffle(len(others), func(i, j int) { others[i], others[j] = others[j], others[i] })
		k := int(float64(len(others)) * connectivity)
		if k < 1 && len(others) > 0 {
			k = 1 // at least one connection
		}
		t.edges[a] = others[:k]
	}
	return t
}

// Neighbors returns the agents that the given agent can send to.
func (t *Topology) Neighbors(agentID string) []string {
	return t.edges[agentID]
}

// CanSend returns whether src can send to dst.
func (t *Topology) CanSend(src, dst string) bool {
	for _, n := range t.edges[src] {
		if n == dst {
			return true
		}
	}
	return false
}

// Agents returns all agents in the topology.
func (t *Topology) Agents() []string {
	out := make([]string, len(t.agents))
	copy(out, t.agents)
	return out
}

// EdgeCount returns the total number of directed edges.
func (t *Topology) EdgeCount() int {
	n := 0
	for _, neighbors := range t.edges {
		n += len(neighbors)
	}
	return n
}

// Connectivity returns the actual connectivity ratio.
func (t *Topology) Connectivity() float64 {
	n := len(t.agents)
	if n <= 1 {
		return 0
	}
	maxEdges := n * (n - 1) // directed
	return float64(t.EdgeCount()) / float64(maxEdges)
}

// AddEdge adds a directed edge from src to dst.
func (t *Topology) AddEdge(src, dst string) {
	if !t.CanSend(src, dst) {
		t.edges[src] = append(t.edges[src], dst)
	}
}

// Edges returns all edges as (src, dst) pairs.
func (t *Topology) Edges() [][2]string {
	var out [][2]string
	// Sort keys for deterministic iteration
	keys := make([]string, 0, len(t.edges))
	for k := range t.edges {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, src := range keys {
		for _, dst := range t.edges[src] {
			out = append(out, [2]string{src, dst})
		}
	}
	return out
}
