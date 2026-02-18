package vllm

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

// ServerState tracks a single vLLM server's status.
type ServerState struct {
	URL            string
	Client         Client
	Models         []string
	GPUCount       int
	MaxTokens      int
	ActiveRequests atomic.Int64
	Healthy        atomic.Bool
	failures       atomic.Int32
}

// Pool manages a set of vLLM servers with health checking and routing.
type Pool struct {
	servers []*ServerState
	router  Router
	mu      sync.RWMutex
	cancel  context.CancelFunc

	healthInterval   time.Duration
	failureThreshold int32
}

type PoolOption func(*Pool)

func WithRouter(r Router) PoolOption {
	return func(p *Pool) { p.router = r }
}

func WithHealthInterval(d time.Duration) PoolOption {
	return func(p *Pool) { p.healthInterval = d }
}

func WithFailureThreshold(n int32) PoolOption {
	return func(p *Pool) { p.failureThreshold = n }
}

// NewPool creates a pool of vLLM servers.
func NewPool(urls []string, opts ...PoolOption) *Pool {
	p := &Pool{
		router:           NewLeastLoadedRouter(),
		healthInterval:   5 * time.Second,
		failureThreshold: 3,
	}
	for _, opt := range opts {
		opt(p)
	}
	for _, url := range urls {
		ss := &ServerState{
			URL:    url,
			Client: NewClient(url),
		}
		ss.Healthy.Store(true) // optimistically healthy
		p.servers = append(p.servers, ss)
	}
	return p
}

// Start begins periodic health checking. Call Stop to shut down.
func (p *Pool) Start(ctx context.Context) {
	ctx, p.cancel = context.WithCancel(ctx)
	go p.healthLoop(ctx)
}

// Stop shuts down the health checker.
func (p *Pool) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
}

// Get selects a server and returns a PoolClient that tracks active requests.
func (p *Pool) Get(ctx context.Context, req *ChatCompletionRequest) (*PoolClient, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	ss, err := p.router.Select(ctx, req, p.servers)
	if err != nil {
		return nil, err
	}
	ss.ActiveRequests.Add(1)
	return &PoolClient{server: ss}, nil
}

// Servers returns a snapshot of all server states.
func (p *Pool) Servers() []*ServerState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]*ServerState, len(p.servers))
	copy(out, p.servers)
	return out
}

func (p *Pool) healthLoop(ctx context.Context) {
	ticker := time.NewTicker(p.healthInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.checkHealth(ctx)
		}
	}
}

func (p *Pool) checkHealth(ctx context.Context) {
	p.mu.RLock()
	servers := make([]*ServerState, len(p.servers))
	copy(servers, p.servers)
	p.mu.RUnlock()

	var wg sync.WaitGroup
	for _, ss := range servers {
		wg.Add(1)
		go func(s *ServerState) {
			defer wg.Done()
			hctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			if err := s.Client.Health(hctx); err != nil {
				n := s.failures.Add(1)
				if n >= p.failureThreshold {
					if s.Healthy.CompareAndSwap(true, false) {
						slog.Warn("vllm server marked unhealthy", "url", s.URL, "failures", n)
					}
				}
			} else {
				s.failures.Store(0)
				if s.Healthy.CompareAndSwap(false, true) {
					slog.Info("vllm server recovered", "url", s.URL)
				}
			}
		}(ss)
	}
	wg.Wait()
}

// PoolClient wraps a server selection. Must call Release when done.
type PoolClient struct {
	server   *ServerState
	released atomic.Bool
}

func (pc *PoolClient) ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	resp, err := pc.server.Client.ChatCompletion(ctx, req)
	if err != nil {
		pc.server.failures.Add(1)
		return nil, err
	}
	return resp, nil
}

func (pc *PoolClient) ChatCompletionStream(ctx context.Context, req *ChatCompletionRequest) (<-chan *ChatCompletionChunk, error) {
	return pc.server.Client.ChatCompletionStream(ctx, req)
}

func (pc *PoolClient) ServerURL() string {
	return pc.server.URL
}

// Release returns the server slot to the pool.
func (pc *PoolClient) Release() {
	if pc.released.CompareAndSwap(false, true) {
		pc.server.ActiveRequests.Add(-1)
	}
}

// MarkUnhealthy records a failure on the underlying server.
func (pc *PoolClient) MarkUnhealthy() {
	pc.server.Healthy.Store(false)
}

// AddServer adds a server to the pool dynamically.
func (p *Pool) AddServer(url string, models []string, gpuCount, maxTokens int) {
	ss := &ServerState{
		URL:       url,
		Client:    NewClient(url),
		Models:    models,
		GPUCount:  gpuCount,
		MaxTokens: maxTokens,
	}
	ss.Healthy.Store(true)
	p.mu.Lock()
	p.servers = append(p.servers, ss)
	p.mu.Unlock()
}

// HealthyCount returns the number of healthy servers.
func (p *Pool) HealthyCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	n := 0
	for _, s := range p.servers {
		if s.Healthy.Load() {
			n++
		}
	}
	return n
}

// ServerCount returns total number of servers.
func (p *Pool) ServerCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.servers)
}

func (p *Pool) String() string {
	return fmt.Sprintf("Pool{servers=%d, healthy=%d}", p.ServerCount(), p.HealthyCount())
}
