package actor

import (
	"fmt"
	"sync"

	hollywoodactor "github.com/anthdm/hollywood/actor"
	"github.com/lingzhi227/go-scheduler/pkg/agent"
)

// AgentFactory creates an Agent from a config.
type AgentFactory func(cfg *agent.AgentConfig) (agent.Agent, error)

// Registry maps agent kinds to their factories and tracks spawned PIDs.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]AgentFactory
	agents    map[string]*hollywoodactor.PID // agentID -> PID
}

func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]AgentFactory),
		agents:    make(map[string]*hollywoodactor.PID),
	}
}

// RegisterFactory registers a factory for a given agent kind.
func (r *Registry) RegisterFactory(kind string, factory AgentFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[kind] = factory
}

// GetFactory returns the factory for a given kind.
func (r *Registry) GetFactory(kind string) (AgentFactory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.factories[kind]
	return f, ok
}

// TrackAgent records an agent PID.
func (r *Registry) TrackAgent(agentID string, pid *hollywoodactor.PID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agents[agentID] = pid
}

// UntrackAgent removes an agent PID.
func (r *Registry) UntrackAgent(agentID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.agents, agentID)
}

// GetAgent returns the PID of a tracked agent.
func (r *Registry) GetAgent(agentID string) (*hollywoodactor.PID, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pid, ok := r.agents[agentID]
	return pid, ok
}

// ListAgents returns a copy of all tracked agent PIDs.
func (r *Registry) ListAgents() map[string]*hollywoodactor.PID {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]*hollywoodactor.PID, len(r.agents))
	for k, v := range r.agents {
		out[k] = v
	}
	return out
}

// CreateAgent uses a factory to create an agent, returns error if kind unknown.
func (r *Registry) CreateAgent(cfg *agent.AgentConfig) (agent.Agent, error) {
	r.mu.RLock()
	factory, ok := r.factories[cfg.Name]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no factory registered for agent kind %q", cfg.Name)
	}
	return factory(cfg)
}
