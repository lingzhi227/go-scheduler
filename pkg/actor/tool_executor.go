package actor

import (
	"context"
	"fmt"
	"log/slog"

	hollywoodactor "github.com/anthdm/hollywood/actor"
	"github.com/lingzhi227/go-scheduler/pkg/agent"
)

// ToolExecutorActor handles tool execution requests from agent actors.
type ToolExecutorActor struct {
	tools map[string]agent.CallableTool
}

// NewToolExecutorProducer creates a Hollywood Producer for a ToolExecutorActor.
func NewToolExecutorProducer(tools map[string]agent.CallableTool) hollywoodactor.Producer {
	return func() hollywoodactor.Receiver {
		return &ToolExecutorActor{tools: tools}
	}
}

func (t *ToolExecutorActor) Receive(ctx *hollywoodactor.Context) {
	switch msg := ctx.Message().(type) {
	case hollywoodactor.Started:
		slog.Info("tool executor started", "tools", len(t.tools))

	case *ToolCallRequest:
		t.handleToolCall(ctx, msg)
	}
}

func (t *ToolExecutorActor) handleToolCall(ctx *hollywoodactor.Context, msg *ToolCallRequest) {
	tool, ok := t.tools[msg.ToolName]
	if !ok {
		response := &ToolCallResponse{
			CallID:   msg.CallID,
			ToolName: msg.ToolName,
			Error:    fmt.Sprintf("unknown tool: %s", msg.ToolName),
		}
		if msg.ReplyTo != nil {
			ctx.Engine().Send(msg.ReplyTo, response)
		}
		return
	}

	result, err := tool.Call(context.Background(), msg.Arguments)
	response := &ToolCallResponse{
		CallID:   msg.CallID,
		ToolName: msg.ToolName,
		Result:   result,
	}
	if err != nil {
		response.Error = err.Error()
		slog.Warn("tool call failed", "tool", msg.ToolName, "error", err)
	}

	if msg.ReplyTo != nil {
		ctx.Engine().Send(msg.ReplyTo, response)
	}
}

// RegisterTool adds a tool to the executor at runtime.
func (t *ToolExecutorActor) RegisterTool(tool agent.CallableTool) {
	t.tools[tool.Declaration().Name] = tool
}
