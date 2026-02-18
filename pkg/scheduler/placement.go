package scheduler

import (
	"github.com/lingzhi227/go-scheduler/pkg/vllm"
)

// PlacementStrategy decides which server to place an agent on.
type PlacementStrategy interface {
	Place(task *ScheduledTask, servers []*vllm.ServerState) *vllm.ServerState
}

// TokenBudgetPlacement selects the server with the most remaining token capacity.
type TokenBudgetPlacement struct{}

func (p *TokenBudgetPlacement) Place(task *ScheduledTask, servers []*vllm.ServerState) *vllm.ServerState {
	var best *vllm.ServerState
	bestCapacity := -1

	for _, s := range servers {
		if !s.Healthy.Load() {
			continue
		}
		// Estimate remaining capacity
		activeReqs := int(s.ActiveRequests.Load())
		capacity := s.MaxTokens - (activeReqs * task.EstimatedTokens)
		if capacity > bestCapacity {
			bestCapacity = capacity
			best = s
		}
	}
	return best
}

// ModelAffinityPlacement selects a server that already has the required model loaded.
type ModelAffinityPlacement struct{}

func (p *ModelAffinityPlacement) Place(task *ScheduledTask, servers []*vllm.ServerState) *vllm.ServerState {
	// Prefer servers with the model already loaded
	for _, s := range servers {
		if !s.Healthy.Load() {
			continue
		}
		for _, m := range s.Models {
			if m == task.Model {
				return s
			}
		}
	}

	// Fallback: any healthy server
	for _, s := range servers {
		if s.Healthy.Load() {
			return s
		}
	}
	return nil
}

// GPUAwarePlacement considers GPU count and load for placement.
// Inspired by Arnold (NeurIPS 2025) topology-aware scheduling.
type GPUAwarePlacement struct{}

func (p *GPUAwarePlacement) Place(task *ScheduledTask, servers []*vllm.ServerState) *vllm.ServerState {
	var best *vllm.ServerState
	bestScore := -1.0

	for _, s := range servers {
		if !s.Healthy.Load() {
			continue
		}
		// Score based on: GPU count, inverse active requests, model affinity
		score := float64(s.GPUCount)
		activeReqs := float64(s.ActiveRequests.Load())
		if activeReqs > 0 {
			score /= activeReqs
		}
		// Bonus for model affinity
		for _, m := range s.Models {
			if m == task.Model {
				score *= 2.0
				break
			}
		}
		if score > bestScore {
			bestScore = score
			best = s
		}
	}
	return best
}
