package backends

import (
	"testing"

	"github.com/lingzhi227/go-scheduler/pkg/memory"
)

func TestInMemorySetGet(t *testing.T) {
	m := NewInMemory()
	if err := m.Set(memory.ScopeWorking, "agent1", "key1", "value1"); err != nil {
		t.Fatal(err)
	}

	val, ok, err := m.Get(memory.ScopeWorking, "agent1", "key1")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}
	if val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}
}

func TestInMemoryDelete(t *testing.T) {
	m := NewInMemory()
	m.Set(memory.ScopeWorking, "a1", "k1", "v1")
	m.Delete(memory.ScopeWorking, "a1", "k1")

	_, ok, _ := m.Get(memory.ScopeWorking, "a1", "k1")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestInMemoryList(t *testing.T) {
	m := NewInMemory()
	m.Set(memory.ScopeWorking, "a1", "key1", "v1")
	m.Set(memory.ScopeWorking, "a1", "key2", "v2")
	m.Set(memory.ScopePlanning, "a1", "other", "v3")

	keys, err := m.List(memory.ScopeWorking, "a1")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestInMemoryScopeIsolation(t *testing.T) {
	m := NewInMemory()
	m.Set(memory.ScopeWorking, "a1", "key", "working")
	m.Set(memory.ScopeProcedural, "a1", "key", "procedural")

	v1, _, _ := m.Get(memory.ScopeWorking, "a1", "key")
	v2, _, _ := m.Get(memory.ScopeProcedural, "a1", "key")

	if v1 == v2 {
		t.Error("scopes should be isolated")
	}
	if v1 != "working" {
		t.Errorf("expected working, got %v", v1)
	}
	if v2 != "procedural" {
		t.Errorf("expected procedural, got %v", v2)
	}
}

func TestInMemoryVectorSearch(t *testing.T) {
	m := NewInMemory()
	m.SetVector(memory.ScopeWorking, "a1", "vec1", []float64{1, 0, 0}, "doc1", nil)
	m.SetVector(memory.ScopeWorking, "a1", "vec2", []float64{0, 1, 0}, "doc2", nil)
	m.SetVector(memory.ScopeWorking, "a1", "vec3", []float64{0.9, 0.1, 0}, "doc3", nil)

	results, err := m.SearchVector(memory.ScopeWorking, "a1", []float64{1, 0, 0}, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// vec1 should be most similar to query [1,0,0]
	if results[0].Key != "vec1" {
		t.Errorf("expected vec1 first, got %s", results[0].Key)
	}
}
