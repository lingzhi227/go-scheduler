package comm

import (
	"log/slog"
	"sync"

	hollywoodactor "github.com/anthdm/hollywood/actor"
	actorpkg "github.com/lingzhi227/go-scheduler/pkg/actor"
)

// ProgressEntry tracks the current status of an agent.
type ProgressEntry struct {
	AgentID        string   `json:"agent_id"`
	TaskID         string   `json:"task_id"`
	Status         string   `json:"status"`
	Turn           int      `json:"turn"`
	Summary        string   `json:"summary"`
	PartialResults []string `json:"partial_results,omitempty"`
}

// ProgressBoard is a Hollywood actor that aggregates agent progress.
// Other agents can query it for peer status.
type ProgressBoard struct {
	mu      sync.RWMutex
	entries map[string]*ProgressEntry // keyed by agentID
}

// NewProgressBoardProducer creates a Hollywood Producer for ProgressBoard.
func NewProgressBoardProducer() hollywoodactor.Producer {
	return func() hollywoodactor.Receiver {
		return &ProgressBoard{
			entries: make(map[string]*ProgressEntry),
		}
	}
}

func (p *ProgressBoard) Receive(ctx *hollywoodactor.Context) {
	switch msg := ctx.Message().(type) {
	case hollywoodactor.Started:
		slog.Info("progress board started")

	case *actorpkg.AgentProgress:
		p.mu.Lock()
		p.entries[msg.AgentID] = &ProgressEntry{
			AgentID:        msg.AgentID,
			TaskID:         msg.TaskID,
			Status:         msg.Status,
			Turn:           msg.Turn,
			Summary:        msg.Summary,
			PartialResults: msg.PartialResults,
		}
		p.mu.Unlock()

	case *actorpkg.GetProgress:
		p.mu.RLock()
		snapshot := make(map[string]*ProgressEntry, len(p.entries))
		for k, v := range p.entries {
			cp := *v
			snapshot[k] = &cp
		}
		p.mu.RUnlock()
		ctx.Respond(&ProgressSnapshot{Entries: snapshot})
	}
}

// ProgressSnapshot is returned by GetProgress.
type ProgressSnapshot struct {
	Entries map[string]*ProgressEntry `json:"entries"`
}

// Snapshot returns a thread-safe copy of all progress entries.
func (p *ProgressBoard) Snapshot() map[string]*ProgressEntry {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make(map[string]*ProgressEntry, len(p.entries))
	for k, v := range p.entries {
		cp := *v
		out[k] = &cp
	}
	return out
}
