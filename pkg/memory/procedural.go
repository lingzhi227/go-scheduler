package memory

import (
	"context"
	"time"
)

// ProceduralEntry represents a learned pattern or procedure (MemGen procedural memory).
type ProceduralEntry struct {
	Key         string         `json:"key"`
	Pattern     string         `json:"pattern"`
	Description string         `json:"description"`
	Examples    []string       `json:"examples,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UsageCount  int            `json:"usage_count"`
}

// ProceduralMemory stores learned patterns and procedures.
type ProceduralMemory struct {
	mem ScopedMemory
}

func NewProceduralMemory(m Memory) *ProceduralMemory {
	return &ProceduralMemory{mem: m.Scoped(ScopeProcedural)}
}

// Store saves a procedural entry.
func (p *ProceduralMemory) Store(ctx context.Context, entry *ProceduralEntry) error {
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}
	return p.mem.Set(ctx, entry.Key, entry)
}

// Recall retrieves a procedural entry by key.
func (p *ProceduralMemory) Recall(ctx context.Context, key string) (*ProceduralEntry, bool, error) {
	val, ok, err := p.mem.Get(ctx, key)
	if err != nil || !ok {
		return nil, false, err
	}
	entry, ok := val.(*ProceduralEntry)
	if !ok {
		return nil, false, nil
	}
	return entry, true, nil
}

// List returns all procedural entry keys.
func (p *ProceduralMemory) List(ctx context.Context) ([]string, error) {
	return p.mem.List(ctx)
}
