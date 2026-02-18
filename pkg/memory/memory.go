package memory

import "context"

// Scope defines the visibility/lifetime of a memory entry.
type Scope string

const (
	ScopeWorking    Scope = "working"    // Current task context (MemGen working memory)
	ScopeProcedural Scope = "procedural" // Learned patterns (MemGen procedural memory)
	ScopePlanning   Scope = "planning"   // Strategy/goals (MemGen planning memory)
	ScopeGlobal     Scope = "global"     // Cross-agent shared state
)

// Memory provides scoped key-value and vector storage.
type Memory interface {
	Set(ctx context.Context, scope Scope, key string, value any) error
	Get(ctx context.Context, scope Scope, key string) (any, bool, error)
	Delete(ctx context.Context, scope Scope, key string) error
	List(ctx context.Context, scope Scope) ([]string, error)
	SearchVector(ctx context.Context, scope Scope, query []float64, limit int) ([]VectorResult, error)
	Scoped(scope Scope) ScopedMemory
}

// ScopedMemory is a view of Memory restricted to a single scope.
type ScopedMemory interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (any, bool, error)
	Delete(ctx context.Context, key string) error
	List(ctx context.Context) ([]string, error)
}

// VectorResult is a single result from a vector similarity search.
type VectorResult struct {
	Key      string         `json:"key"`
	Score    float64        `json:"score"`
	Value    any            `json:"value"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// memoryImpl implements Memory using a Backend.
type memoryImpl struct {
	backend Backend
	ownerID string // agent or session ID
}

// New creates a new Memory backed by the given Backend.
func New(backend Backend, ownerID string) Memory {
	return &memoryImpl{backend: backend, ownerID: ownerID}
}

func (m *memoryImpl) Set(ctx context.Context, scope Scope, key string, value any) error {
	return m.backend.Set(scope, m.ownerID, key, value)
}

func (m *memoryImpl) Get(ctx context.Context, scope Scope, key string) (any, bool, error) {
	return m.backend.Get(scope, m.ownerID, key)
}

func (m *memoryImpl) Delete(ctx context.Context, scope Scope, key string) error {
	return m.backend.Delete(scope, m.ownerID, key)
}

func (m *memoryImpl) List(ctx context.Context, scope Scope) ([]string, error) {
	return m.backend.List(scope, m.ownerID)
}

func (m *memoryImpl) SearchVector(ctx context.Context, scope Scope, query []float64, limit int) ([]VectorResult, error) {
	return m.backend.SearchVector(scope, m.ownerID, query, limit)
}

func (m *memoryImpl) Scoped(scope Scope) ScopedMemory {
	return &scopedMemory{m: m, scope: scope}
}

type scopedMemory struct {
	m     *memoryImpl
	scope Scope
}

func (s *scopedMemory) Set(ctx context.Context, key string, value any) error {
	return s.m.Set(ctx, s.scope, key, value)
}

func (s *scopedMemory) Get(ctx context.Context, key string) (any, bool, error) {
	return s.m.Get(ctx, s.scope, key)
}

func (s *scopedMemory) Delete(ctx context.Context, key string) error {
	return s.m.Delete(ctx, s.scope, key)
}

func (s *scopedMemory) List(ctx context.Context) ([]string, error) {
	return s.m.List(ctx, s.scope)
}
