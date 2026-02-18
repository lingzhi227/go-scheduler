package orchestrate

import (
	"fmt"

	"github.com/lingzhi227/go-scheduler/pkg/agent"
)

// BuildParallelGraph creates a graph where multiple agents work in parallel on
// the same task, then results are joined.
func BuildParallelGraph(id, goal string, agentConfigs []*agent.AgentConfig) *TaskGraph {
	g := NewTaskGraph(id, goal)

	start := &Node{ID: "start", Type: NodeTypeStart, Label: "Start"}
	join := &Node{ID: "join", Type: NodeTypeJoin, Label: "Aggregate results"}
	end := &Node{ID: "end", Type: NodeTypeEnd, Label: "End"}
	g.AddNode(start)
	g.AddNode(join)
	g.AddNode(end)

	for i, cfg := range agentConfigs {
		nodeID := fmt.Sprintf("agent_%d", i)
		n := &Node{
			ID:          nodeID,
			Type:        NodeTypeAgent,
			Label:       cfg.Name,
			AgentConfig: cfg,
			Input:       map[string]any{"goal": goal},
		}
		g.AddNode(n)
		g.AddEdge("start", nodeID)
		g.AddEdge(nodeID, "join")
	}
	g.AddEdge("join", "end")

	return g
}

// BuildChainGraph creates a sequential chain of agents.
// Each agent receives the output of the previous one.
// Paper: Chain-of-Agents (NeurIPS 2024).
func BuildChainGraph(id, goal string, agentConfigs []*agent.AgentConfig) *TaskGraph {
	g := NewTaskGraph(id, goal)

	start := &Node{ID: "start", Type: NodeTypeStart, Label: "Start"}
	end := &Node{ID: "end", Type: NodeTypeEnd, Label: "End"}
	g.AddNode(start)
	g.AddNode(end)

	prev := "start"
	for i, cfg := range agentConfigs {
		nodeID := fmt.Sprintf("agent_%d", i)
		n := &Node{
			ID:          nodeID,
			Type:        NodeTypeAgent,
			Label:       cfg.Name,
			AgentConfig: cfg,
			Input:       map[string]any{"goal": goal},
		}
		g.AddNode(n)
		g.AddEdge(prev, nodeID)
		prev = nodeID
	}
	g.AddEdge(prev, "end")

	return g
}

// BuildDebateGraph creates a multi-round debate between agents.
// In each round, agents see each other's previous outputs.
func BuildDebateGraph(id, goal string, agentConfigs []*agent.AgentConfig, rounds int) *TaskGraph {
	g := NewTaskGraph(id, goal)

	start := &Node{ID: "start", Type: NodeTypeStart, Label: "Start"}
	end := &Node{ID: "end", Type: NodeTypeEnd, Label: "End"}
	g.AddNode(start)
	g.AddNode(end)

	prev := "start"
	for r := 0; r < rounds; r++ {
		var roundNodes []string
		for i, cfg := range agentConfigs {
			nodeID := fmt.Sprintf("round%d_agent%d", r, i)
			n := &Node{
				ID:          nodeID,
				Type:        NodeTypeAgent,
				Label:       fmt.Sprintf("Round %d: %s", r+1, cfg.Name),
				AgentConfig: cfg,
				Input: map[string]any{
					"goal":  goal,
					"round": r + 1,
				},
			}
			g.AddNode(n)
			g.AddEdge(prev, nodeID)
			roundNodes = append(roundNodes, nodeID)
		}

		// Add join after each round
		joinID := fmt.Sprintf("round%d_join", r)
		join := &Node{ID: joinID, Type: NodeTypeJoin, Label: fmt.Sprintf("Round %d summary", r+1)}
		g.AddNode(join)
		for _, nid := range roundNodes {
			g.AddEdge(nid, joinID)
		}
		prev = joinID
	}

	g.AddEdge(prev, "end")
	return g
}

// BuildMapReduceGraph creates a map-reduce pattern: fan out to workers, then reduce.
func BuildMapReduceGraph(id string, subGoals []string, workerConfig, reducerConfig *agent.AgentConfig) *TaskGraph {
	g := NewTaskGraph(id, "map-reduce")

	start := &Node{ID: "start", Type: NodeTypeStart, Label: "Start"}
	join := &Node{ID: "join", Type: NodeTypeJoin, Label: "Reduce"}
	end := &Node{ID: "end", Type: NodeTypeEnd, Label: "End"}
	g.AddNode(start)
	g.AddNode(join)
	g.AddNode(end)

	for i, sg := range subGoals {
		nodeID := fmt.Sprintf("worker_%d", i)
		n := &Node{
			ID:          nodeID,
			Type:        NodeTypeAgent,
			Label:       fmt.Sprintf("Worker: %s", sg),
			AgentConfig: workerConfig,
			Input:       map[string]any{"goal": sg},
		}
		g.AddNode(n)
		g.AddEdge("start", nodeID)
		g.AddEdge(nodeID, "join")
	}

	// Reducer agent synthesizes joined results
	reducer := &Node{
		ID:          "reducer",
		Type:        NodeTypeAgent,
		Label:       "Reducer",
		AgentConfig: reducerConfig,
		Input:       map[string]any{"goal": "Synthesize the results from all workers"},
	}
	g.AddNode(reducer)
	g.AddEdge("join", "reducer")
	g.AddEdge("reducer", "end")

	return g
}
