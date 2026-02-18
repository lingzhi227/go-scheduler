package vllm

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"
)

// Router selects a server from the pool for a given request.
type Router interface {
	Select(ctx context.Context, req *ChatCompletionRequest, servers []*ServerState) (*ServerState, error)
}

// RoundRobinRouter distributes requests evenly across healthy servers.
type RoundRobinRouter struct {
	counter atomic.Uint64
}

func NewRoundRobinRouter() *RoundRobinRouter {
	return &RoundRobinRouter{}
}

func (r *RoundRobinRouter) Select(_ context.Context, _ *ChatCompletionRequest, servers []*ServerState) (*ServerState, error) {
	healthy := filterHealthy(servers)
	if len(healthy) == 0 {
		return nil, fmt.Errorf("no healthy servers available")
	}
	idx := r.counter.Add(1) % uint64(len(healthy))
	return healthy[idx], nil
}

// LeastLoadedRouter picks the server with the fewest active requests.
type LeastLoadedRouter struct{}

func NewLeastLoadedRouter() *LeastLoadedRouter {
	return &LeastLoadedRouter{}
}

func (r *LeastLoadedRouter) Select(_ context.Context, _ *ChatCompletionRequest, servers []*ServerState) (*ServerState, error) {
	healthy := filterHealthy(servers)
	if len(healthy) == 0 {
		return nil, fmt.Errorf("no healthy servers available")
	}
	var best *ServerState
	bestLoad := int64(math.MaxInt64)
	for _, s := range healthy {
		load := s.ActiveRequests.Load()
		if load < bestLoad {
			bestLoad = load
			best = s
		}
	}
	return best, nil
}

// ModelAffinityRouter routes to servers that have the requested model loaded.
// Falls back to least-loaded if no affinity match.
type ModelAffinityRouter struct {
	fallback *LeastLoadedRouter
}

func NewModelAffinityRouter() *ModelAffinityRouter {
	return &ModelAffinityRouter{fallback: NewLeastLoadedRouter()}
}

func (r *ModelAffinityRouter) Select(ctx context.Context, req *ChatCompletionRequest, servers []*ServerState) (*ServerState, error) {
	healthy := filterHealthy(servers)
	if len(healthy) == 0 {
		return nil, fmt.Errorf("no healthy servers available")
	}

	// Filter for model affinity
	var affine []*ServerState
	for _, s := range healthy {
		for _, m := range s.Models {
			if m == req.Model {
				affine = append(affine, s)
				break
			}
		}
	}

	if len(affine) > 0 {
		return r.fallback.Select(ctx, req, affine)
	}
	return r.fallback.Select(ctx, req, healthy)
}

func filterHealthy(servers []*ServerState) []*ServerState {
	var result []*ServerState
	for _, s := range servers {
		if s.Healthy.Load() {
			result = append(result, s)
		}
	}
	return result
}
