# go-scheduler — Distributed Agent Operating System

<p align="center">
  <strong>Orchestrate dozens-to-hundreds of LLM agents working in parallel toward collective goals.</strong>
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge" alt="MIT License"></a>
  <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.22+"></a>
</p>

**go-scheduler** is a Go scheduler that orchestrates open-source agentic LLMs on GPU clusters (4 GPUs/node, vLLM, batch=10, 10-80K tokens). It handles tool calling, shared memory, inter-agent communication, DAG orchestration, and progress visibility. Think: distributed Agent OS.

## Architecture

```
                    REST API (:8080)
                        |
                   +---------+
                   | Executor |  (DAG orchestration)
                   +---------+
                   /    |     \
            +-----+ +-----+ +-----+
            |Agent| |Agent| |Agent|   (Hollywood actors)
            +-----+ +-----+ +-----+
               |        |       |
            +--+--------+-------+--+
            |    vLLM Pool         |   (health checks, routing)
            +--+--------+-------+--+
               |        |       |
            [GPU 0]  [GPU 1]  [GPU 2]  (vLLM servers)
```

## Features

- **vLLM Pool** — connection pool with health checks, round-robin / least-loaded / model-affinity routing
- **ReAct Agent Loop** — LLM call -> tool calls -> repeat, with streaming events
- **Actor Engine** — Hollywood actors for agents, tool executors, and supervision trees (OneForOne / AllForOne / RestForOne)
- **Multi-Level Memory** — working / procedural / planning tiers (MemGen paper), scoped per agent or shared
- **Shared Spaces** — named memory spaces with ACL (reader / writer / admin) for cross-agent state
- **Inter-Agent Comms** — sparse / ring / star / full topologies, P2P channels, broadcast, ProgressBoard
- **DAG Orchestration** — task graphs with topological sort, parallel execution, LLM-driven goal decomposition
- **Aggregation Strategies** — majority vote (NeurIPS 2025), chain-of-agents (NeurIPS 2024), LLM merge
- **Pre-built Patterns** — parallel, chain, debate, map-reduce graphs
- **GPU-Aware Scheduling** — priority queue, token budget / model affinity / GPU-aware placement
- **Durable Execution** — workflow engine with checkpointing for agent state recovery
- **Observability** — Prometheus metrics (`go_scheduler_agents_active`, `go_scheduler_llm_latency_seconds`, `go_scheduler_tokens_total`) + OpenTelemetry tracing
- **Cluster Mode** — Hollywood cluster with pluggable discovery providers

## Quick Start

Runtime: **Go 1.22+**

```bash
git clone https://github.com/lingzhi227/go-scheduler.git
cd go-scheduler

go build ./cmd/go-scheduler

# Run with default config (expects a vLLM server on localhost:8000)
./go-scheduler

# Run with config file
./go-scheduler -config config.json

# Run with environment overrides
GO_SCHEDULER_VLLM_URLS=http://gpu1:8000,http://gpu2:8000 GO_SCHEDULER_ADDR=:9090 ./go-scheduler
```

## Configuration

JSON config file or environment variables:

```json
{
  "server": { "addr": ":8080" },
  "vllm": {
    "router": "least_loaded",
    "servers": [
      { "url": "http://gpu1:8000", "models": ["llama-3.1-70b"], "gpu_count": 4, "max_tokens": 80000 },
      { "url": "http://gpu2:8000", "models": ["qwen-2.5-72b"], "gpu_count": 4, "max_tokens": 80000 }
    ]
  },
  "agents": [
    { "id": "researcher", "name": "general", "system_prompt": "You are a research agent.", "max_turns": 15 },
    { "id": "coder", "name": "general", "system_prompt": "You are a coding agent.", "max_turns": 10 }
  ],
  "memory": { "backend": "inmemory" }
}
```

| Env var | Description |
|---------|-------------|
| `GO_SCHEDULER_ADDR` | REST API listen address (default `:8080`) |
| `GO_SCHEDULER_VLLM_URLS` | Comma-separated vLLM server URLs |
| `GO_SCHEDULER_MEMORY_BACKEND` | `inmemory` or `redis` |
| `GO_SCHEDULER_REDIS_ADDR` | Redis address for distributed memory |

## REST API

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/tasks` | Submit a task to an agent |
| `GET` | `/tasks/{id}` | Get task status and result |
| `GET` | `/tasks` | List all tasks |
| `GET` | `/agents` | List active agents |
| `GET` | `/progress` | Get all agents' progress |
| `POST` | `/goals` | Submit a high-level goal for decomposition |
| `GET` | `/health` | Health check |

### Submit a task

```bash
curl -X POST http://localhost:8080/tasks \
  -H 'Content-Type: application/json' \
  -d '{"goal": "Summarize the key findings in this dataset", "agent_id": "researcher"}'
```

### Check task status

```bash
curl http://localhost:8080/tasks/<task-id>
```

## Project Structure

```
go-scheduler/
├── cmd/go-scheduler/main.go         # Entry point
├── api/rest/server.go               # REST API
├── internal/config/config.go        # JSON/env config loading
├── pkg/
│   ├── vllm/                        # vLLM client, pool, routing
│   │   ├── types.go                 #   ChatCompletion request/response types
│   │   ├── client.go                #   OpenAI-compatible HTTP client
│   │   ├── pool.go                  #   Connection pool + health checks
│   │   └── router.go                #   Round-robin, least-loaded, model-affinity
│   ├── agent/                       # Core agent abstractions
│   │   ├── agent.go                 #   Agent interface + BaseAgent
│   │   ├── tool.go                  #   CallableTool interface + registry
│   │   ├── event.go                 #   Event types (thinking, llm_call, tool_call, done, error)
│   │   └── llm.go                   #   ReAct loop implementation
│   ├── actor/                       # Hollywood actor layer
│   │   ├── agent_actor.go           #   AgentActor wrapping agent.Agent
│   │   ├── tool_executor.go         #   ToolExecutor actor pool
│   │   ├── supervisor.go            #   Supervision trees (OneForOne/AllForOne/RestForOne)
│   │   ├── messages.go              #   RunTask, TaskCompleted, ToolCallRequest, etc.
│   │   └── registry.go              #   Agent kind -> factory registry
│   ├── memory/                      # Multi-level memory (MemGen-inspired)
│   │   ├── memory.go                #   Memory interface with scopes
│   │   ├── backend.go               #   Backend interface
│   │   ├── working.go               #   Working memory (current session)
│   │   ├── procedural.go            #   Procedural memory (learned patterns)
│   │   ├── planning.go              #   Planning memory (strategy/goals)
│   │   ├── shared.go                #   SharedSpace with ACL
│   │   └── backends/
│   │       ├── inmemory.go          #   In-process map + vector search
│   │       └── redis.go             #   Redis-backed distributed store
│   ├── comm/                        # Inter-agent communication
│   │   ├── topology.go              #   Sparse/ring/star/full topologies
│   │   ├── channel.go               #   P2P messaging over actors
│   │   ├── broadcast.go             #   Multicast within topology
│   │   └── progress.go              #   ProgressBoard actor
│   ├── orchestrate/                 # DAG orchestration
│   │   ├── graph.go                 #   TaskGraph with topological sort
│   │   ├── executor.go              #   Parallel executor (level-by-level)
│   │   ├── decompose.go             #   LLM-driven goal -> TaskGraph
│   │   ├── aggregate.go             #   Vote / chain / LLM-merge strategies
│   │   └── patterns.go              #   Pre-built: parallel, chain, debate, map-reduce
│   ├── scheduler/                   # GPU-aware scheduling
│   │   ├── scheduler.go             #   Priority queue + vLLM server selection
│   │   └── placement.go             #   Token budget / model affinity / GPU-aware placement
│   ├── durable/                     # Durable execution
│   │   ├── workflow.go              #   Workflow engine interface
│   │   └── checkpoint.go            #   Agent state checkpointing
│   ├── observe/                     # Observability
│   │   ├── metrics.go               #   Prometheus metrics
│   │   └── trace.go                 #   OpenTelemetry tracing
│   └── cluster/                     # Cluster management
│       ├── cluster.go               #   Hollywood cluster wrapper
│       └── provider.go              #   Static discovery provider
└── github/                          # Reference projects (not part of build)
```

## Key Dependencies

| Dependency | Purpose |
|------------|---------|
| [hollywood](https://github.com/anthdm/hollywood) | Actor engine (10M msg/sec, built-in cluster) |
| [go-redis](https://github.com/redis/go-redis) | Distributed memory backend |
| [uuid](https://github.com/google/uuid) | Task/graph ID generation |
| [prometheus/client_golang](https://github.com/prometheus/client_golang) | Metrics |
| [opentelemetry](https://opentelemetry.io/docs/languages/go/) | Distributed tracing |

## Paper References

| Paper | Conference | Applied in |
|-------|-----------|------------|
| Sparse Communication Topology | EMNLP 2024 | `pkg/comm/topology.go` — 30-50% connectivity default |
| Vote over Debate | NeurIPS 2025 | `pkg/orchestrate/aggregate.go` — majority voting |
| Chain-of-Agents | NeurIPS 2024 | `pkg/orchestrate/patterns.go` — sequential chain with communication units |
| MemGen | ICLR 2026 | `pkg/memory/` — working / procedural / planning tiers |
| Arnold | NeurIPS 2025 | `pkg/scheduler/placement.go` — topology-aware GPU placement |

## Development

```bash
# Build
go build ./...

# Run tests
go test ./...

# Vet
go vet ./...

# Run with race detector
go test -race ./...
```

## Verification Checklist

- **Phase 1**: Start 2+ vLLM servers (or mock), submit a task via REST, see agent complete it with tool calls
- **Phase 2**: Run 2 agents, verify Agent B can read Agent A's progress via ProgressBoard; verify shared memory writes are visible cross-agent
- **Phase 3**: Submit a complex goal, verify decomposition into DAG, parallel execution, and aggregated result
- **Phase 4**: Kill a vLLM server mid-task, verify circuit breaker kicks in and agent retries on healthy server; restart an agent, verify it recovers from checkpoint

## License

MIT
