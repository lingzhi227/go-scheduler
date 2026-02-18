package memory

import (
	"context"
	"time"
)

// PlanEntry represents a planning/goal entry (MemGen planning memory).
type PlanEntry struct {
	Key        string    `json:"key"`
	Goal       string    `json:"goal"`
	Strategy   string    `json:"strategy"`
	SubGoals   []string  `json:"sub_goals,omitempty"`
	Status     string    `json:"status"` // pending, active, completed, abandoned
	Priority   int       `json:"priority"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// PlanningMemory stores strategy and goal information.
type PlanningMemory struct {
	mem ScopedMemory
}

func NewPlanningMemory(m Memory) *PlanningMemory {
	return &PlanningMemory{mem: m.Scoped(ScopePlanning)}
}

// SetGoal stores or updates a planning goal.
func (p *PlanningMemory) SetGoal(ctx context.Context, entry *PlanEntry) error {
	now := time.Now()
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = now
	}
	entry.UpdatedAt = now
	return p.mem.Set(ctx, entry.Key, entry)
}

// GetGoal retrieves a planning goal by key.
func (p *PlanningMemory) GetGoal(ctx context.Context, key string) (*PlanEntry, bool, error) {
	val, ok, err := p.mem.Get(ctx, key)
	if err != nil || !ok {
		return nil, false, err
	}
	entry, ok := val.(*PlanEntry)
	if !ok {
		return nil, false, nil
	}
	return entry, true, nil
}

// ListGoals returns all planning keys.
func (p *PlanningMemory) ListGoals(ctx context.Context) ([]string, error) {
	return p.mem.List(ctx)
}
