package observe

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for AgentSched.
var Metrics = struct {
	AgentsActive      prometheus.Gauge
	TasksSubmitted    prometheus.Counter
	TasksCompleted    *prometheus.CounterVec
	LLMCallsTotal    prometheus.Counter
	LLMLatency        prometheus.Histogram
	TokensTotal       *prometheus.CounterVec
	ToolCallsTotal    *prometheus.CounterVec
	VLLMServerHealth  *prometheus.GaugeVec
	SchedulerQueueLen prometheus.Gauge
}{
	AgentsActive: promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "agentsched",
		Name:      "agents_active",
		Help:      "Number of currently active agents.",
	}),
	TasksSubmitted: promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "agentsched",
		Name:      "tasks_submitted_total",
		Help:      "Total number of tasks submitted.",
	}),
	TasksCompleted: promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "agentsched",
		Name:      "tasks_completed_total",
		Help:      "Total tasks completed by status.",
	}, []string{"status"}), // success, error
	LLMCallsTotal: promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "agentsched",
		Name:      "llm_calls_total",
		Help:      "Total number of LLM API calls.",
	}),
	LLMLatency: promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "agentsched",
		Name:      "llm_latency_seconds",
		Help:      "LLM call latency in seconds.",
		Buckets:   []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60},
	}),
	TokensTotal: promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "agentsched",
		Name:      "tokens_total",
		Help:      "Total tokens consumed.",
	}, []string{"type"}), // prompt, completion
	ToolCallsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "agentsched",
		Name:      "tool_calls_total",
		Help:      "Total tool invocations.",
	}, []string{"tool", "status"}), // tool name, success/error
	VLLMServerHealth: promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "agentsched",
		Name:      "vllm_server_healthy",
		Help:      "Health status of vLLM servers (1=healthy, 0=unhealthy).",
	}, []string{"url"}),
	SchedulerQueueLen: promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "agentsched",
		Name:      "scheduler_queue_length",
		Help:      "Number of tasks in the scheduler queue.",
	}),
}
