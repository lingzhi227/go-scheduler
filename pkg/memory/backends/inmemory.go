package backends

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/lingzhi227/go-scheduler/pkg/memory"
)

// InMemory implements memory.Backend with in-process maps.
type InMemory struct {
	mu      sync.RWMutex
	data    map[string]map[string]any     // scope:scopeID -> key -> value
	vectors map[string]map[string]vecEntry // scope:scopeID -> key -> vector
}

type vecEntry struct {
	embedding []float64
	metadata  map[string]any
	value     any
}

func NewInMemory() *InMemory {
	return &InMemory{
		data:    make(map[string]map[string]any),
		vectors: make(map[string]map[string]vecEntry),
	}
}

func scopeKey(scope memory.Scope, scopeID string) string {
	return string(scope) + ":" + scopeID
}

func (m *InMemory) Set(scope memory.Scope, scopeID, key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	sk := scopeKey(scope, scopeID)
	if m.data[sk] == nil {
		m.data[sk] = make(map[string]any)
	}
	m.data[sk][key] = value
	return nil
}

func (m *InMemory) Get(scope memory.Scope, scopeID, key string) (any, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sk := scopeKey(scope, scopeID)
	if m.data[sk] == nil {
		return nil, false, nil
	}
	v, ok := m.data[sk][key]
	return v, ok, nil
}

func (m *InMemory) Delete(scope memory.Scope, scopeID, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	sk := scopeKey(scope, scopeID)
	if m.data[sk] != nil {
		delete(m.data[sk], key)
	}
	return nil
}

func (m *InMemory) List(scope memory.Scope, scopeID string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sk := scopeKey(scope, scopeID)
	if m.data[sk] == nil {
		return nil, nil
	}
	keys := make([]string, 0, len(m.data[sk]))
	for k := range m.data[sk] {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys, nil
}

func (m *InMemory) SearchVector(scope memory.Scope, scopeID string, query []float64, limit int) ([]memory.VectorResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sk := scopeKey(scope, scopeID)
	vecs := m.vectors[sk]
	if len(vecs) == 0 {
		return nil, nil
	}

	type scored struct {
		key   string
		score float64
		entry vecEntry
	}
	var results []scored
	for k, v := range vecs {
		score := cosineSimilarity(query, v.embedding)
		results = append(results, scored{key: k, score: score, entry: v})
	}
	sort.Slice(results, func(i, j int) bool { return results[i].score > results[j].score })

	if limit > len(results) {
		limit = len(results)
	}
	out := make([]memory.VectorResult, limit)
	for i := 0; i < limit; i++ {
		out[i] = memory.VectorResult{
			Key:      results[i].key,
			Score:    results[i].score,
			Value:    results[i].entry.value,
			Metadata: results[i].entry.metadata,
		}
	}
	return out, nil
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	denom := math.Sqrt(normA) * math.Sqrt(normB)
	if denom == 0 {
		return 0
	}
	return dot / denom
}

// SetVector stores a vector embedding.
func (m *InMemory) SetVector(scope memory.Scope, scopeID, key string, embedding []float64, value any, metadata map[string]any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	sk := scopeKey(scope, scopeID)
	if m.vectors[sk] == nil {
		m.vectors[sk] = make(map[string]vecEntry)
	}
	m.vectors[sk][key] = vecEntry{
		embedding: embedding,
		metadata:  metadata,
		value:     value,
	}
	return nil
}

func (m *InMemory) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return fmt.Sprintf("InMemory{scopes=%d}", len(m.data))
}
