package observe

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "go-scheduler"

// Tracer returns the OpenTelemetry tracer for go-scheduler.
func Tracer() trace.Tracer {
	return otel.Tracer(tracerName)
}

// StartAgentTurnSpan starts a span for an agent's turn.
func StartAgentTurnSpan(ctx context.Context, agentID, taskID string, turn int) (context.Context, trace.Span) {
	return Tracer().Start(ctx, "agent.turn",
		trace.WithAttributes(
			attribute.String("agent.id", agentID),
			attribute.String("task.id", taskID),
			attribute.Int("turn", turn),
		),
	)
}

// StartLLMCallSpan starts a span for an LLM call.
func StartLLMCallSpan(ctx context.Context, model string, messageCount int) (context.Context, trace.Span) {
	return Tracer().Start(ctx, "llm.call",
		trace.WithAttributes(
			attribute.String("llm.model", model),
			attribute.Int("llm.message_count", messageCount),
		),
	)
}

// StartToolCallSpan starts a span for a tool invocation.
func StartToolCallSpan(ctx context.Context, toolName, callID string) (context.Context, trace.Span) {
	return Tracer().Start(ctx, "tool.call",
		trace.WithAttributes(
			attribute.String("tool.name", toolName),
			attribute.String("tool.call_id", callID),
		),
	)
}

// StartGraphExecutionSpan starts a span for task graph execution.
func StartGraphExecutionSpan(ctx context.Context, graphID, goal string) (context.Context, trace.Span) {
	return Tracer().Start(ctx, "graph.execute",
		trace.WithAttributes(
			attribute.String("graph.id", graphID),
			attribute.String("graph.goal", goal),
		),
	)
}
