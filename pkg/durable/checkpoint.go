package durable

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/lingzhi227/go-scheduler/pkg/vllm"
)

// Checkpoint captures an agent's state for recovery.
type Checkpoint struct {
	AgentID   string         `json:"agent_id"`
	TaskID    string         `json:"task_id"`
	Turn      int            `json:"turn"`
	Messages  []vllm.Message `json:"messages"`
	State     map[string]any `json:"state,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// CheckpointStore persists and retrieves checkpoints.
type CheckpointStore interface {
	Save(cp *Checkpoint) error
	Load(agentID, taskID string) (*Checkpoint, error)
	Delete(agentID, taskID string) error
	List(agentID string) ([]*Checkpoint, error)
}

// InMemoryCheckpointStore stores checkpoints in memory.
type InMemoryCheckpointStore struct {
	mu          sync.RWMutex
	checkpoints map[string]*Checkpoint // key: agentID:taskID
}

func NewInMemoryCheckpointStore() *InMemoryCheckpointStore {
	return &InMemoryCheckpointStore{
		checkpoints: make(map[string]*Checkpoint),
	}
}

func cpKey(agentID, taskID string) string {
	return agentID + ":" + taskID
}

func (s *InMemoryCheckpointStore) Save(cp *Checkpoint) error {
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = time.Now()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checkpoints[cpKey(cp.AgentID, cp.TaskID)] = cp
	return nil
}

func (s *InMemoryCheckpointStore) Load(agentID, taskID string) (*Checkpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cp, ok := s.checkpoints[cpKey(agentID, taskID)]
	if !ok {
		return nil, fmt.Errorf("checkpoint not found: %s:%s", agentID, taskID)
	}
	return cp, nil
}

func (s *InMemoryCheckpointStore) Delete(agentID, taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.checkpoints, cpKey(agentID, taskID))
	return nil
}

func (s *InMemoryCheckpointStore) List(agentID string) ([]*Checkpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []*Checkpoint
	for k, v := range s.checkpoints {
		if len(k) > len(agentID) && k[:len(agentID)] == agentID {
			result = append(result, v)
		}
	}
	return result, nil
}

// Serialize serializes a checkpoint to JSON bytes.
func (cp *Checkpoint) Serialize() ([]byte, error) {
	return json.Marshal(cp)
}

// DeserializeCheckpoint deserializes a checkpoint from JSON bytes.
func DeserializeCheckpoint(data []byte) (*Checkpoint, error) {
	var cp Checkpoint
	if err := json.Unmarshal(data, &cp); err != nil {
		return nil, err
	}
	return &cp, nil
}
