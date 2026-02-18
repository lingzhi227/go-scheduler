# Repository Guidelines

- Repo: `github.com/lingzhi227/go-scheduler`
- Language: Go (1.22+). Module path: `github.com/lingzhi227/go-scheduler`.

## Project Structure

- Source code lives under `cmd/`, `pkg/`, `internal/`, `api/`.
- Tests: colocated `*_test.go` files. Use `_test` package suffix to avoid import cycles when testing against subpackages.
- Reference projects: `github/` directory contains cloned repos for API reference only. They are fenced off by their own `go.mod` and are **not** part of the build. Never edit files under `github/`.
- Config: `internal/config/` loads JSON config with env var overrides (`AGENTSCHED_*` prefix).
- REST API: `api/rest/` uses Go 1.22+ `http.ServeMux` routing (`METHOD /path`).
- Entry point: `cmd/agentsched/main.go`.

## Package Map

| Package | Responsibility | Key interfaces |
|---------|---------------|----------------|
| `pkg/vllm` | vLLM server communication | `Client`, `Pool`, `Router` |
| `pkg/agent` | Core agent abstractions | `Agent`, `CallableTool`, `Event` |
| `pkg/actor` | Hollywood actor wrappers | `AgentActor`, `SupervisorActor`, `ToolExecutorActor` |
| `pkg/memory` | Multi-level scoped memory | `Memory`, `Backend`, `SharedSpace` |
| `pkg/comm` | Inter-agent communication | `Topology`, `Channel`, `ProgressBoard` |
| `pkg/orchestrate` | DAG task graphs | `TaskGraph`, `Executor`, `Decomposer`, `Aggregator` |
| `pkg/scheduler` | Priority queue + placement | `Scheduler`, `PlacementStrategy` |
| `pkg/durable` | Workflow engine + checkpoints | `WorkflowEngine`, `CheckpointStore` |
| `pkg/observe` | Prometheus + OpenTelemetry | `Metrics`, `Tracer()` |
| `pkg/cluster` | Multi-node cluster | `Cluster`, `Provider` |

## Build, Test, and Lint Commands

```bash
# Build all packages
go build ./...

# Run all tests
go test ./...

# Vet
go vet ./...

# Run tests with race detector
go test -race ./...

# Build binary
go build -o agentsched ./cmd/agentsched

# Tidy dependencies
go mod tidy
```

## Coding Style

- **Language**: Go. Follow standard Go conventions (`gofmt`, `go vet`).
- **Naming**: exported types use `PascalCase`, unexported use `camelCase`. Package names are lowercase single words.
- **Errors**: return `error` as the last return value. Wrap with `fmt.Errorf("context: %w", err)`.
- **Concurrency**: prefer channels and `sync` primitives. Actor message passing via Hollywood for inter-agent communication.
- **Logging**: use `log/slog` (structured logging). No `fmt.Println` in library code.
- **Interfaces**: define in the consuming package when possible. Keep interfaces small (1-3 methods).
- **Tests**: use standard `testing` package. No test frameworks. Table-driven tests preferred.
- **Files**: keep under ~500 LOC. Split when it improves clarity.
- **Comments**: add for non-obvious logic. Exported types and functions should have doc comments.
- **Imports**: group stdlib, external, internal (separated by blank lines).

## Architecture Decisions

- **Actor engine**: Hollywood (not protoactor/ergo). Agents are IO-bound on LLM calls; Hollywood's simple `Receiver` interface wins over heavier alternatives.
- **Message types**: defined in `pkg/actor/messages.go`. All actor communication uses typed Go structs (not protobuf) for simplicity.
- **vLLM communication**: OpenAI-compatible HTTP API. No gRPC. Types in `pkg/vllm/types.go` mirror the OpenAI chat completion spec.
- **Memory scopes**: four tiers — `working` (session), `procedural` (patterns), `planning` (goals), `global` (shared). Scopes are string constants, not iota.
- **Topology**: default sparse connectivity (30-50%) per EMNLP 2024 sparse communication paper. All topology types in `pkg/comm/topology.go`.
- **Aggregation**: default majority vote per NeurIPS 2025. LLM-merge available for richer synthesis.
- **Config**: JSON file + `AGENTSCHED_*` env vars. No YAML dependency (keep deps minimal).

## Key Patterns

### Agent ReAct Loop (`pkg/agent/llm.go`)

```
Build messages (system + task + context)
Loop (up to maxTurns):
  1. Call vLLM ChatCompletion
  2. If tool_calls in response → execute tools, append results, continue
  3. If stop/no tool calls → emit Done event, return
```

### Actor Message Flow

```
REST API → RunTask message → AgentActor
  AgentActor runs ReAct loop in goroutine
  AgentActor sends AgentProgress → ProgressBoard
  AgentActor sends TaskCompleted/TaskFailed → collector actor → REST response
```

### Task Graph Execution (`pkg/orchestrate/executor.go`)

```
TopologicalSort → levels (groups of independent nodes)
For each level:
  Execute all nodes in parallel (goroutines)
  Wait for level completion
  Store outputs in graph state for downstream nodes
Join nodes aggregate predecessor outputs via Aggregator
```

### Supervisor Restart (`pkg/actor/supervisor.go`)

```
ChildFailed message received →
  Check restart budget (max 5 restarts per minute window)
  OneForOne: restart only the failed child
  AllForOne: restart all children
  RestForOne: restart failed child + all children started after it
```

## Dependency Notes

- **Hollywood** (`github.com/anthdm/hollywood`): actor engine. Core types: `Engine`, `Receiver`, `Context`, `PID`, `Producer`. Lifecycle messages: `Initialized`, `Started`, `Stopped`.
- **go-redis** (`github.com/redis/go-redis/v9`): used only in `pkg/memory/backends/redis.go`. Not required for dev (in-memory backend is default).
- **Prometheus**: metrics registered via `promauto` in `pkg/observe/metrics.go`. Counters and gauges under `agentsched_*` namespace.
- **OpenTelemetry**: tracing spans in `pkg/observe/trace.go`. Tracer name: `agentsched`.

## Testing

- Tests live alongside source: `foo.go` → `foo_test.go`.
- Use `*_test` (external test) package when importing subpackages to avoid import cycles (see `pkg/memory/shared_test.go`).
- Mock vLLM servers: use `httptest.NewServer` returning canned `ChatCompletionResponse` (see `pkg/vllm/pool_test.go`).
- No mocking frameworks. Interfaces + test doubles.
- Run `go test ./...` before committing.

## Git & Commits

- Commit messages: concise, action-oriented (e.g., `vllm: add model affinity router`).
- Scope prefix matches package: `vllm:`, `agent:`, `actor:`, `memory:`, `comm:`, `orchestrate:`, `scheduler:`, `rest:`, `config:`.
- Group related changes; avoid bundling unrelated refactors.
- Do not commit `github/` reference project changes.

## Common Tasks

### Adding a new tool

1. Implement `agent.CallableTool` interface (see `pkg/agent/tool.go`)
2. Register in the tool registry or pass directly to an agent
3. Tools are automatically converted to vLLM function definitions via `agent.ToVLLMTools()`

### Adding a new agent type

1. Implement `agent.Agent` interface
2. Register factory in `pkg/actor/registry.go` via `RegisterFactory(kind, factory)`
3. Add to config YAML/JSON under `agents[]`

### Adding a new topology

1. Add constructor in `pkg/comm/topology.go` (e.g., `NewTreeTopology`)
2. Returns `*Topology` with edges populated
3. Register in comm setup code

### Adding a new aggregation strategy

1. Implement `orchestrate.Aggregator` interface (single method: `Aggregate(ctx, outputs) (string, error)`)
2. Use in executor or pass to `NewExecutor()`

### Adding a new memory backend

1. Implement `memory.Backend` interface in `pkg/memory/backends/`
2. Wire in `cmd/agentsched/main.go` config switch

### Adding a new graph pattern

1. Add builder function in `pkg/orchestrate/patterns.go`
2. Returns `*TaskGraph` with nodes and edges pre-wired

## Agent-Specific Notes

- When adding new packages, follow the existing pattern: interface in the package root, implementations in subpackages or the same file.
- The `github/` directory is read-only reference material. Key references:
  - `github/hollywood/actor/` — Hollywood API (Engine, Context, Receiver)
  - `github/trpc-agent-go/agent/` — Agent/tool interface patterns
  - `github/go-agent/` — SharedSpace, SubAgentTool, ToolCatalog
  - `github/agentfield/sdk/go/agent/` — MemoryBackend, scoped memory
  - `github/go-workflows/workflow/` — Durable workflow patterns
  - `github/ergo/act/` — Supervisor strategy patterns
- When modifying actor message types, update both `pkg/actor/messages.go` and all receivers that handle those messages.
- Pool client wrappers (`poolClientWrapper`) adapt `*vllm.PoolClient` to `vllm.Client` interface. Both `pkg/actor/agent_actor.go` and `pkg/orchestrate/executor.go` have their own copies — keep them in sync or extract to a shared helper.
- The REST API uses Go 1.22+ enhanced `ServeMux` routing patterns (`METHOD /path`). No external router dependency.
- Hollywood actor lifecycle: handle `actor.Started` for init, `actor.Stopped` for cleanup. Never block in `Receive` — spawn goroutines for long work.
- Memory import cycle: `pkg/memory/backends/` imports `pkg/memory` for types. Tests in `pkg/memory/` that need backends must use the `memory_test` package.
