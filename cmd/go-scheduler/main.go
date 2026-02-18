package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	hollywoodactor "github.com/anthdm/hollywood/actor"
	actorpkg "github.com/lingzhi227/go-scheduler/pkg/actor"
	"github.com/lingzhi227/go-scheduler/pkg/agent"
	"github.com/lingzhi227/go-scheduler/pkg/comm"
	"github.com/lingzhi227/go-scheduler/pkg/memory/backends"
	"github.com/lingzhi227/go-scheduler/pkg/orchestrate"
	"github.com/lingzhi227/go-scheduler/pkg/vllm"

	"github.com/lingzhi227/go-scheduler/api/rest"
	"github.com/lingzhi227/go-scheduler/internal/config"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	// Setup structured logging
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	// Load config
	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Create vLLM pool
	var urls []string
	for _, s := range cfg.VLLM.Servers {
		urls = append(urls, s.URL)
	}
	if len(urls) == 0 {
		slog.Warn("no vLLM servers configured, using mock endpoint")
		urls = []string{"http://localhost:8000"}
	}

	var router vllm.Router
	switch cfg.VLLM.Router {
	case "round_robin":
		router = vllm.NewRoundRobinRouter()
	case "model_affinity":
		router = vllm.NewModelAffinityRouter()
	default:
		router = vllm.NewLeastLoadedRouter()
	}

	pool := vllm.NewPool(urls, vllm.WithRouter(router))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool.Start(ctx)
	defer pool.Stop()

	// Add server metadata
	for i, s := range cfg.VLLM.Servers {
		if i < pool.ServerCount() {
			servers := pool.Servers()
			servers[i].Models = s.Models
			servers[i].GPUCount = s.GPUCount
			servers[i].MaxTokens = s.MaxTokens
		}
	}

	// Create Hollywood actor engine
	engine, err := hollywoodactor.NewEngine(hollywoodactor.NewEngineConfig())
	if err != nil {
		slog.Error("failed to create actor engine", "error", err)
		os.Exit(1)
	}

	// Create memory backend
	memBackend := backends.NewInMemory()
	_ = memBackend

	// Create agent registry
	registry := actorpkg.NewRegistry()

	// Register a default general-purpose agent factory
	registry.RegisterFactory("general", func(cfg *agent.AgentConfig) (agent.Agent, error) {
		return &generalAgent{
			BaseAgent: agent.NewBaseAgent(cfg.ID, cfg.SystemPrompt, cfg.Model, nil, cfg.MaxTurns),
		}, nil
	})

	// Spawn progress board
	progressPID := engine.Spawn(comm.NewProgressBoardProducer(), "progress_board")

	// Build child specs from config
	var childSpecs []actorpkg.ChildSpec
	for _, agentCfg := range cfg.Agents {
		acfg := agentCfg // capture
		ag, err := registry.CreateAgent(&acfg)
		if err != nil {
			slog.Error("failed to create agent", "id", acfg.ID, "error", err)
			continue
		}
		childSpecs = append(childSpecs, actorpkg.ChildSpec{
			Kind:     acfg.ID,
			Producer: actorpkg.NewAgentActorProducer(ag, pool, nil, progressPID),
		})
	}

	// If no agents configured, create a default one
	if len(childSpecs) == 0 {
		defaultCfg := &agent.AgentConfig{
			ID:           "default",
			Name:         "general",
			SystemPrompt: "You are a helpful assistant.",
			MaxTurns:     10,
		}
		ag, _ := registry.CreateAgent(defaultCfg)
		childSpecs = append(childSpecs, actorpkg.ChildSpec{
			Kind:     "default",
			Producer: actorpkg.NewAgentActorProducer(ag, pool, nil, progressPID),
		})
	}

	// Spawn supervisor with OneForOne strategy
	supervisorPID := engine.Spawn(
		actorpkg.NewSupervisorProducer(actorpkg.OneForOne, childSpecs, registry),
		"supervisor",
	)

	// Create orchestrator
	executor := orchestrate.NewExecutor(engine, pool, registry, &orchestrate.VoteAggregator{})

	// Create and start REST server
	apiServer := rest.NewServer(engine, registry, supervisorPID, progressPID, executor)

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		slog.Info("shutting down...")
		apiServer.Shutdown(ctx)
		engine.Poison(supervisorPID)
		engine.Poison(progressPID)
		cancel()
	}()

	slog.Info("go-scheduler starting",
		"addr", cfg.Server.Addr,
		"vllm_servers", len(urls),
		"agents", len(childSpecs),
	)

	if err := apiServer.Start(cfg.Server.Addr); err != nil && err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

// generalAgent is a simple default agent implementation.
type generalAgent struct {
	*agent.BaseAgent
}

func (a *generalAgent) Run(ctx context.Context, task *agent.Task) (<-chan *agent.Event, error) {
	tools := make(map[string]agent.CallableTool)
	for _, t := range a.Tools() {
		tools[t.Declaration().Name] = t
	}
	// The actual loop will be handled by the actor via RunAgentLoop
	ch := make(chan *agent.Event, 1)
	close(ch)
	return ch, nil
}
