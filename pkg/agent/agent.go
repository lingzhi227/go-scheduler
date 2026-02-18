package agent

import "context"

// Agent is the core abstraction for an LLM-powered agent.
type Agent interface {
	// ID returns a unique identifier for this agent.
	ID() string

	// Run starts the agent on a task, returning an event stream.
	Run(ctx context.Context, task *Task) (<-chan *Event, error)

	// Tools returns the tools available to this agent.
	Tools() []CallableTool

	// SystemPrompt returns the system prompt for this agent.
	SystemPrompt() string
}

// Task represents a unit of work for an agent.
type Task struct {
	ID          string            `json:"id"`
	Goal        string            `json:"goal"`
	Context     string            `json:"context,omitempty"`
	Input       map[string]any    `json:"input,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	MaxTurns    int               `json:"max_turns,omitempty"`
	ParentID    string            `json:"parent_id,omitempty"`
	Priority    int               `json:"priority,omitempty"`
}

// TaskResult holds the final output of a completed task.
type TaskResult struct {
	TaskID  string         `json:"task_id"`
	AgentID string         `json:"agent_id"`
	Output  string         `json:"output"`
	Data    map[string]any `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
	Usage   TaskUsage      `json:"usage"`
}

// TaskUsage tracks resource usage for a task.
type TaskUsage struct {
	TotalTokens      int           `json:"total_tokens"`
	PromptTokens     int           `json:"prompt_tokens"`
	CompletionTokens int           `json:"completion_tokens"`
	LLMCalls         int           `json:"llm_calls"`
	ToolCalls        int           `json:"tool_calls"`
	Turns            int           `json:"turns"`
}

// AgentConfig describes how to instantiate an agent.
type AgentConfig struct {
	ID           string            `json:"id" yaml:"id"`
	Name         string            `json:"name" yaml:"name"`
	Model        string            `json:"model" yaml:"model"`
	SystemPrompt string            `json:"system_prompt" yaml:"system_prompt"`
	Tools        []string          `json:"tools" yaml:"tools"`
	MaxTurns     int               `json:"max_turns" yaml:"max_turns"`
	Temperature  *float64          `json:"temperature,omitempty" yaml:"temperature,omitempty"`
	Labels       map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

// BaseAgent provides a minimal Agent implementation.
type BaseAgent struct {
	id           string
	systemPrompt string
	tools        []CallableTool
	maxTurns     int
	model        string
}

func NewBaseAgent(id, systemPrompt, model string, tools []CallableTool, maxTurns int) *BaseAgent {
	if maxTurns <= 0 {
		maxTurns = 10
	}
	return &BaseAgent{
		id:           id,
		systemPrompt: systemPrompt,
		tools:        tools,
		maxTurns:     maxTurns,
		model:        model,
	}
}

func (a *BaseAgent) ID() string             { return a.id }
func (a *BaseAgent) SystemPrompt() string    { return a.systemPrompt }
func (a *BaseAgent) Tools() []CallableTool   { return a.tools }
func (a *BaseAgent) MaxTurns() int           { return a.maxTurns }
func (a *BaseAgent) Model() string           { return a.model }
