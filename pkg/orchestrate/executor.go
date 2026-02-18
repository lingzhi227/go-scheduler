package orchestrate

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	hollywoodactor "github.com/anthdm/hollywood/actor"
	"github.com/lingzhi227/go-scheduler/pkg/actor"
	"github.com/lingzhi227/go-scheduler/pkg/agent"
	"github.com/lingzhi227/go-scheduler/pkg/vllm"
)

// ExecutionStatus tracks the overall execution state.
type ExecutionStatus string

const (
	ExecPending  ExecutionStatus = "pending"
	ExecRunning  ExecutionStatus = "running"
	ExecDone     ExecutionStatus = "done"
	ExecFailed   ExecutionStatus = "failed"
)

// NodeResult holds the result of executing a single node.
type NodeResult struct {
	NodeID    string         `json:"node_id"`
	Status    NodeStatus     `json:"status"`
	Output    string         `json:"output"`
	Data      map[string]any `json:"data,omitempty"`
	Error     string         `json:"error,omitempty"`
	StartedAt time.Time     `json:"started_at"`
	DoneAt    time.Time      `json:"done_at"`
}

// ExecutionResult holds the complete execution result.
type ExecutionResult struct {
	GraphID   string                 `json:"graph_id"`
	Status    ExecutionStatus        `json:"status"`
	Results   map[string]*NodeResult `json:"results"`
	Output    string                 `json:"output"`
	StartedAt time.Time             `json:"started_at"`
	DoneAt    time.Time              `json:"done_at"`
}

// Executor runs a TaskGraph using the actor engine.
type Executor struct {
	engine      *hollywoodactor.Engine
	pool        *vllm.Pool
	registry    *actor.Registry
	aggregator  Aggregator
}

func NewExecutor(engine *hollywoodactor.Engine, pool *vllm.Pool, registry *actor.Registry, agg Aggregator) *Executor {
	if agg == nil {
		agg = &VoteAggregator{}
	}
	return &Executor{
		engine:     engine,
		pool:       pool,
		registry:   registry,
		aggregator: agg,
	}
}

// Execute runs the task graph, processing nodes level by level.
func (e *Executor) Execute(ctx context.Context, graph *TaskGraph) (*ExecutionResult, error) {
	levels, err := graph.TopologicalSort()
	if err != nil {
		return nil, fmt.Errorf("topological sort: %w", err)
	}

	result := &ExecutionResult{
		GraphID:   graph.ID,
		Status:    ExecRunning,
		Results:   make(map[string]*NodeResult),
		StartedAt: time.Now(),
	}

	for _, level := range levels {
		if ctx.Err() != nil {
			result.Status = ExecFailed
			return result, ctx.Err()
		}

		// Execute all nodes in this level in parallel
		var wg sync.WaitGroup
		var mu sync.Mutex
		levelFailed := false

		for _, nodeID := range level {
			node := graph.Nodes[nodeID]
			if node == nil {
				continue
			}

			// Skip start/end nodes
			if node.Type == NodeTypeStart || node.Type == NodeTypeEnd {
				mu.Lock()
				result.Results[nodeID] = &NodeResult{
					NodeID:    nodeID,
					Status:    NodeDone,
					StartedAt: time.Now(),
					DoneAt:    time.Now(),
				}
				mu.Unlock()
				continue
			}

			// Skip join nodes (they aggregate results)
			if node.Type == NodeTypeJoin {
				preds := graph.Predecessors(nodeID)
				var outputs []string
				mu.Lock()
				for _, p := range preds {
					if r, ok := result.Results[p]; ok && r.Status == NodeDone {
						outputs = append(outputs, r.Output)
					}
				}
				mu.Unlock()

				aggregated, aggErr := e.aggregator.Aggregate(ctx, outputs)
				nr := &NodeResult{
					NodeID:    nodeID,
					StartedAt: time.Now(),
					DoneAt:    time.Now(),
				}
				if aggErr != nil {
					nr.Status = NodeFailed
					nr.Error = aggErr.Error()
					levelFailed = true
				} else {
					nr.Status = NodeDone
					nr.Output = aggregated
				}
				mu.Lock()
				result.Results[nodeID] = nr
				mu.Unlock()
				continue
			}

			wg.Add(1)
			go func(n *Node) {
				defer wg.Done()
				nr := e.executeNode(ctx, n, graph)
				mu.Lock()
				result.Results[n.ID] = nr
				if nr.Status == NodeFailed {
					levelFailed = true
				}
				// Store output in graph state for downstream nodes
				if nr.Output != "" {
					graph.SetState(n.ID+".output", nr.Output)
				}
				mu.Unlock()
			}(node)
		}
		wg.Wait()

		if levelFailed {
			result.Status = ExecFailed
			result.DoneAt = time.Now()
			return result, fmt.Errorf("execution failed at level containing nodes %v", level)
		}
	}

	result.Status = ExecDone
	result.DoneAt = time.Now()

	// Set final output from the last non-end node
	for i := len(levels) - 1; i >= 0; i-- {
		for _, nodeID := range levels[i] {
			if r, ok := result.Results[nodeID]; ok && r.Output != "" {
				result.Output = r.Output
				return result, nil
			}
		}
	}

	return result, nil
}

func (e *Executor) executeNode(ctx context.Context, node *Node, graph *TaskGraph) *NodeResult {
	nr := &NodeResult{
		NodeID:    node.ID,
		Status:    NodeRunning,
		StartedAt: time.Now(),
	}

	switch node.Type {
	case NodeTypeAgent:
		nr = e.executeAgentNode(ctx, node, graph)
	case NodeTypeTool:
		nr = e.executeToolNode(ctx, node, graph)
	case NodeTypeDecision:
		nr = e.executeDecisionNode(ctx, node, graph)
	default:
		nr.Status = NodeFailed
		nr.Error = fmt.Sprintf("unknown node type: %s", node.Type)
	}

	nr.DoneAt = time.Now()
	return nr
}

func (e *Executor) executeAgentNode(ctx context.Context, node *Node, graph *TaskGraph) *NodeResult {
	nr := &NodeResult{NodeID: node.ID, StartedAt: time.Now()}

	if node.AgentConfig == nil {
		nr.Status = NodeFailed
		nr.Error = "no agent config for agent node"
		return nr
	}

	// Create agent from config
	ag, err := e.registry.CreateAgent(node.AgentConfig)
	if err != nil {
		nr.Status = NodeFailed
		nr.Error = fmt.Sprintf("create agent: %v", err)
		return nr
	}

	// Build goal from node input + predecessor outputs
	goal := node.Label
	if input, ok := node.Input["goal"]; ok {
		goal = fmt.Sprintf("%v", input)
	}

	// Gather context from predecessors
	preds := graph.Predecessors(node.ID)
	var predContext string
	for _, p := range preds {
		if output, ok := graph.GetState(p + ".output"); ok {
			predContext += fmt.Sprintf("Result from %s: %v\n", p, output)
		}
	}

	task := &agent.Task{
		ID:      graph.ID + ":" + node.ID,
		Goal:    goal,
		Context: predContext,
	}
	if mt, ok := node.Input["max_turns"]; ok {
		if turns, ok := mt.(int); ok {
			task.MaxTurns = turns
		}
	}

	// Get a vLLM client and run the agent loop
	poolClient, err := e.pool.Get(ctx, &vllm.ChatCompletionRequest{Model: node.AgentConfig.Model})
	if err != nil {
		nr.Status = NodeFailed
		nr.Error = fmt.Sprintf("get vllm client: %v", err)
		return nr
	}
	defer poolClient.Release()

	// Build tool map
	tools := make(map[string]agent.CallableTool)
	for _, t := range ag.Tools() {
		tools[t.Declaration().Name] = t
	}

	wc := &poolClientWrapper{pc: poolClient}
	events := agent.RunAgentLoop(ctx, wc, ag, task, tools)

	for evt := range events {
		switch evt.Type {
		case agent.EventTypeDone:
			nr.Status = NodeDone
			nr.Output = evt.Message
			if evt.Result != nil {
				nr.Output = evt.Result.Output
			}
		case agent.EventTypeError:
			nr.Status = NodeFailed
			nr.Error = evt.Error
		}
	}

	if nr.Status == "" {
		nr.Status = NodeFailed
		nr.Error = "agent loop ended without result"
	}
	return nr
}

func (e *Executor) executeToolNode(ctx context.Context, node *Node, _ *TaskGraph) *NodeResult {
	nr := &NodeResult{NodeID: node.ID, StartedAt: time.Now()}
	nr.Status = NodeFailed
	nr.Error = "tool node execution not yet implemented"
	slog.Warn("tool node skipped", "node", node.ID, "tool", node.ToolName)
	return nr
}

func (e *Executor) executeDecisionNode(ctx context.Context, node *Node, _ *TaskGraph) *NodeResult {
	nr := &NodeResult{NodeID: node.ID, StartedAt: time.Now()}
	nr.Status = NodeDone
	nr.Output = "default"
	return nr
}

// poolClientWrapper adapts vllm.PoolClient to vllm.Client interface.
type poolClientWrapper struct {
	pc *vllm.PoolClient
}

func (w *poolClientWrapper) ChatCompletion(ctx context.Context, req *vllm.ChatCompletionRequest) (*vllm.ChatCompletionResponse, error) {
	return w.pc.ChatCompletion(ctx, req)
}

func (w *poolClientWrapper) ChatCompletionStream(ctx context.Context, req *vllm.ChatCompletionRequest) (<-chan *vllm.ChatCompletionChunk, error) {
	return w.pc.ChatCompletionStream(ctx, req)
}

func (w *poolClientWrapper) Health(ctx context.Context) error { return nil }
func (w *poolClientWrapper) BaseURL() string                  { return w.pc.ServerURL() }
