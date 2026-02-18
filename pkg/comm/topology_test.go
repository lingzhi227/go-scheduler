package comm

import (
	"testing"
)

func TestFullTopology(t *testing.T) {
	agents := []string{"a", "b", "c"}
	topo := NewFullTopology(agents)

	if !topo.CanSend("a", "b") {
		t.Error("a should be able to send to b")
	}
	if !topo.CanSend("b", "a") {
		t.Error("b should be able to send to a")
	}
	if topo.CanSend("a", "a") {
		t.Error("a should not send to itself")
	}

	expected := 3 * 2 // 3 agents * 2 connections each
	if topo.EdgeCount() != expected {
		t.Errorf("expected %d edges, got %d", expected, topo.EdgeCount())
	}
}

func TestRingTopology(t *testing.T) {
	agents := []string{"a", "b", "c"}
	topo := NewRingTopology(agents)

	if !topo.CanSend("a", "b") {
		t.Error("a should send to b")
	}
	if !topo.CanSend("b", "c") {
		t.Error("b should send to c")
	}
	if !topo.CanSend("c", "a") {
		t.Error("c should send to a (ring)")
	}
	if topo.CanSend("a", "c") {
		t.Error("a should not send directly to c in ring")
	}

	if topo.EdgeCount() != 3 {
		t.Errorf("expected 3 edges, got %d", topo.EdgeCount())
	}
}

func TestStarTopology(t *testing.T) {
	topo := NewStarTopology("hub", []string{"s1", "s2", "s3"})

	if !topo.CanSend("hub", "s1") {
		t.Error("hub should send to s1")
	}
	if !topo.CanSend("s2", "hub") {
		t.Error("s2 should send to hub")
	}
	if topo.CanSend("s1", "s2") {
		t.Error("spokes should not send to each other")
	}

	// 3 from hub + 3 to hub = 6
	if topo.EdgeCount() != 6 {
		t.Errorf("expected 6 edges, got %d", topo.EdgeCount())
	}
}

func TestSparseTopology(t *testing.T) {
	agents := make([]string, 10)
	for i := range agents {
		agents[i] = string(rune('a' + i))
	}
	topo := NewSparseTopology(agents, 0.3)

	// Each agent should have at least 1 connection
	for _, a := range agents {
		neighbors := topo.Neighbors(a)
		if len(neighbors) < 1 {
			t.Errorf("agent %s has no neighbors", a)
		}
	}

	// Connectivity should be approximately 30%
	conn := topo.Connectivity()
	if conn < 0.1 || conn > 0.6 {
		t.Errorf("connectivity %.2f is outside expected range [0.1, 0.6]", conn)
	}
}

func TestSparseTopologySingleAgent(t *testing.T) {
	topo := NewSparseTopology([]string{"a"}, 0.5)
	if topo.EdgeCount() != 0 {
		t.Error("single agent should have no edges")
	}
}
