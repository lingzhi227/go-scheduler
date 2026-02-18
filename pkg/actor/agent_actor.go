package actor

import (
	"context"
	"log/slog"

	hollywoodactor "github.com/anthdm/hollywood/actor"
	"github.com/lingzhi227/go-scheduler/pkg/agent"
	"github.com/lingzhi227/go-scheduler/pkg/vllm"
)

// AgentActor wraps an agent.Agent as a Hollywood Receiver.
type AgentActor struct {
	agent         agent.Agent
	pool          *vllm.Pool
	toolExecutor  *hollywoodactor.PID
	progressBoard *hollywoodactor.PID
	cancel        context.CancelFunc
}

// NewAgentActorProducer returns a Hollywood Producer for an AgentActor.
func NewAgentActorProducer(ag agent.Agent, pool *vllm.Pool, toolExec, progress *hollywoodactor.PID) hollywoodactor.Producer {
	return func() hollywoodactor.Receiver {
		return &AgentActor{
			agent:         ag,
			pool:          pool,
			toolExecutor:  toolExec,
			progressBoard: progress,
		}
	}
}

func (a *AgentActor) Receive(ctx *hollywoodactor.Context) {
	switch msg := ctx.Message().(type) {
	case hollywoodactor.Started:
		slog.Info("agent actor started", "agent", a.agent.ID())

	case hollywoodactor.Stopped:
		slog.Info("agent actor stopped", "agent", a.agent.ID())
		if a.cancel != nil {
			a.cancel()
		}

	case *RunTask:
		a.handleRunTask(ctx, msg)
	}
}

func (a *AgentActor) handleRunTask(ctx *hollywoodactor.Context, msg *RunTask) {
	task := msg.Task
	slog.Info("agent starting task", "agent", a.agent.ID(), "task", task.ID, "goal", task.Goal)

	// Build tool map
	tools := make(map[string]agent.CallableTool)
	for _, t := range a.agent.Tools() {
		tools[t.Declaration().Name] = t
	}

	// Get a vLLM client from pool
	poolClient, err := a.pool.Get(ctx.Context(), &vllm.ChatCompletionRequest{})
	if err != nil {
		slog.Error("failed to get vllm client", "error", err)
		if msg.ReplyTo != nil {
			ctx.Engine().Send(msg.ReplyTo, &TaskFailed{
				TaskID:  task.ID,
				AgentID: a.agent.ID(),
				Error:   err.Error(),
			})
		}
		return
	}

	// Run the agent loop in a goroutine
	runCtx, cancel := context.WithCancel(ctx.Context())
	a.cancel = cancel

	go func() {
		defer poolClient.Release()
		defer cancel()

		// Create a simple wrapper client for the pool client
		wc := &poolClientWrapper{pc: poolClient}
		events := agent.RunAgentLoop(runCtx, wc, a.agent, task, tools)

		for evt := range events {
			// Forward progress to board
			if a.progressBoard != nil {
				ctx.Engine().Send(a.progressBoard, &AgentProgress{
					AgentID: evt.AgentID,
					TaskID:  evt.TaskID,
					Status:  string(evt.Type),
					Turn:    evt.Turn,
					Summary: evt.Message,
				})
			}

			// On completion/error, send result back
			switch evt.Type {
			case agent.EventTypeDone:
				if msg.ReplyTo != nil {
					ctx.Engine().Send(msg.ReplyTo, &TaskCompleted{Result: evt.Result})
				}
			case agent.EventTypeError:
				if msg.ReplyTo != nil {
					ctx.Engine().Send(msg.ReplyTo, &TaskFailed{
						TaskID:  task.ID,
						AgentID: a.agent.ID(),
						Error:   evt.Error,
					})
				}
			}
		}
	}()
}

// poolClientWrapper adapts PoolClient to the vllm.Client interface.
type poolClientWrapper struct {
	pc *vllm.PoolClient
}

func (w *poolClientWrapper) ChatCompletion(ctx context.Context, req *vllm.ChatCompletionRequest) (*vllm.ChatCompletionResponse, error) {
	return w.pc.ChatCompletion(ctx, req)
}

func (w *poolClientWrapper) ChatCompletionStream(ctx context.Context, req *vllm.ChatCompletionRequest) (<-chan *vllm.ChatCompletionChunk, error) {
	return w.pc.ChatCompletionStream(ctx, req)
}

func (w *poolClientWrapper) Health(ctx context.Context) error {
	return nil // pool handles health
}

func (w *poolClientWrapper) BaseURL() string {
	return w.pc.ServerURL()
}
