package orchestrate

import (
	"testing"

	"github.com/lingzhi227/go-scheduler/pkg/agent"
)

func TestTaskGraphTopologicalSort(t *testing.T) {
	g := NewTaskGraph("test", "test goal")
	g.AddNode(&Node{ID: "start", Type: NodeTypeStart})
	g.AddNode(&Node{ID: "a", Type: NodeTypeAgent})
	g.AddNode(&Node{ID: "b", Type: NodeTypeAgent})
	g.AddNode(&Node{ID: "c", Type: NodeTypeAgent})
	g.AddNode(&Node{ID: "end", Type: NodeTypeEnd})

	g.AddEdge("start", "a")
	g.AddEdge("start", "b")
	g.AddEdge("a", "c")
	g.AddEdge("b", "c")
	g.AddEdge("c", "end")

	levels, err := g.TopologicalSort()
	if err != nil {
		t.Fatal(err)
	}

	if len(levels) != 4 {
		t.Fatalf("expected 4 levels, got %d", len(levels))
	}

	// Level 0: start
	if len(levels[0]) != 1 || levels[0][0] != "start" {
		t.Errorf("level 0 should be [start], got %v", levels[0])
	}

	// Level 1: a and b (parallel)
	if len(levels[1]) != 2 {
		t.Errorf("level 1 should have 2 nodes (a,b), got %d", len(levels[1]))
	}

	// Level 2: c
	if len(levels[2]) != 1 || levels[2][0] != "c" {
		t.Errorf("level 2 should be [c], got %v", levels[2])
	}
}

func TestTaskGraphCycleDetection(t *testing.T) {
	g := NewTaskGraph("test", "cycle test")
	g.AddNode(&Node{ID: "a", Type: NodeTypeAgent})
	g.AddNode(&Node{ID: "b", Type: NodeTypeAgent})
	g.AddEdge("a", "b")
	g.AddEdge("b", "a")

	_, err := g.TopologicalSort()
	if err == nil {
		t.Error("expected cycle detection error")
	}
}

func TestBuildParallelGraph(t *testing.T) {
	configs := []*agent.AgentConfig{
		{ID: "a1", Name: "agent1"},
		{ID: "a2", Name: "agent2"},
		{ID: "a3", Name: "agent3"},
	}
	g := BuildParallelGraph("test", "parallel goal", configs)

	levels, err := g.TopologicalSort()
	if err != nil {
		t.Fatal(err)
	}

	// Should be: start -> [agent_0, agent_1, agent_2] -> join -> end
	if len(levels) != 4 {
		t.Fatalf("expected 4 levels, got %d: %v", len(levels), levels)
	}
	if len(levels[1]) != 3 {
		t.Errorf("expected 3 parallel agents, got %d", len(levels[1]))
	}
}

func TestBuildChainGraph(t *testing.T) {
	configs := []*agent.AgentConfig{
		{ID: "a1", Name: "agent1"},
		{ID: "a2", Name: "agent2"},
	}
	g := BuildChainGraph("test", "chain goal", configs)

	levels, err := g.TopologicalSort()
	if err != nil {
		t.Fatal(err)
	}

	// start -> agent_0 -> agent_1 -> end
	if len(levels) != 4 {
		t.Fatalf("expected 4 levels, got %d", len(levels))
	}
}

func TestGraphState(t *testing.T) {
	g := NewTaskGraph("test", "state test")
	g.SetState("key1", "value1")

	val, ok := g.GetState("key1")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}

	_, ok = g.GetState("nonexistent")
	if ok {
		t.Error("expected nonexistent key to not exist")
	}
}
