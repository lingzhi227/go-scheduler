package orchestrate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/lingzhi227/go-scheduler/pkg/agent"
	"github.com/lingzhi227/go-scheduler/pkg/vllm"
)

// Decomposer uses an LLM to break a high-level goal into a TaskGraph.
type Decomposer struct {
	llm vllm.Client
}

func NewDecomposer(llm vllm.Client) *Decomposer {
	return &Decomposer{llm: llm}
}

const decompositionPrompt = `You are a task decomposition engine. Given a high-level goal and a list of available agents, break the goal into a directed acyclic graph (DAG) of subtasks.

Available agents:
%s

Respond with a JSON object:
{
  "nodes": [
    {"id": "unique_id", "type": "agent", "label": "what this step does", "agent_name": "which agent to use", "input": {"goal": "specific subtask description"}}
  ],
  "edges": [
    {"from": "node_id_1", "to": "node_id_2"}
  ]
}

Node types: "agent" (delegates to an agent), "join" (aggregates results from predecessors).
Always include a "start" node and an "end" node.
Use edges to define dependencies. Nodes with no dependencies from each other can run in parallel.
Respond ONLY with the JSON, no other text.`

// AgentInfo describes an available agent for decomposition.
type AgentInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Tools       []string `json:"tools,omitempty"`
}

// decompositionResponse is the expected LLM output format.
type decompositionResponse struct {
	Nodes []struct {
		ID        string         `json:"id"`
		Type      string         `json:"type"`
		Label     string         `json:"label"`
		AgentName string         `json:"agent_name,omitempty"`
		Input     map[string]any `json:"input,omitempty"`
	} `json:"nodes"`
	Edges []struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"edges"`
}

// Decompose takes a goal and available agents, returns a TaskGraph.
func (d *Decomposer) Decompose(ctx context.Context, goal string, agents []AgentInfo) (*TaskGraph, error) {
	agentsJSON, _ := json.MarshalIndent(agents, "", "  ")

	req := &vllm.ChatCompletionRequest{
		Messages: []vllm.Message{
			{Role: "system", Content: fmt.Sprintf(decompositionPrompt, string(agentsJSON))},
			{Role: "user", Content: goal},
		},
	}

	resp, err := d.llm.ChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("llm decomposition call: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in decomposition response")
	}

	var dr decompositionResponse
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &dr); err != nil {
		return nil, fmt.Errorf("parse decomposition: %w (raw: %s)", err, resp.Choices[0].Message.Content)
	}

	graphID := uuid.New().String()
	graph := NewTaskGraph(graphID, goal)

	for _, n := range dr.Nodes {
		node := &Node{
			ID:    n.ID,
			Type:  NodeType(n.Type),
			Label: n.Label,
			Input: n.Input,
		}
		if n.AgentName != "" {
			node.AgentConfig = &agent.AgentConfig{
				Name: n.AgentName,
			}
		}
		if err := graph.AddNode(node); err != nil {
			return nil, fmt.Errorf("add node %s: %w", n.ID, err)
		}
	}

	for _, e := range dr.Edges {
		if err := graph.AddEdge(e.From, e.To); err != nil {
			return nil, fmt.Errorf("add edge %s->%s: %w", e.From, e.To, err)
		}
	}

	return graph, nil
}
