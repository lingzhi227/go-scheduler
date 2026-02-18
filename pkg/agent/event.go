package agent

import "time"

// EventType enumerates the kinds of events an agent can emit.
type EventType string

const (
	EventTypeThinking   EventType = "thinking"    // Agent is reasoning
	EventTypeLLMCall    EventType = "llm_call"     // LLM request started
	EventTypeLLMResult  EventType = "llm_result"   // LLM response received
	EventTypeToolCall   EventType = "tool_call"    // Tool invocation started
	EventTypeToolResult EventType = "tool_result"  // Tool result received
	EventTypeMessage    EventType = "message"      // Intermediate message
	EventTypeProgress   EventType = "progress"     // Progress update
	EventTypeDone       EventType = "done"         // Task completed
	EventTypeError      EventType = "error"        // Error occurred
)

// Event represents a single event in an agent's execution.
type Event struct {
	Type      EventType      `json:"type"`
	AgentID   string         `json:"agent_id"`
	TaskID    string         `json:"task_id"`
	Turn      int            `json:"turn"`
	Timestamp time.Time      `json:"timestamp"`
	Data      any            `json:"data,omitempty"`
	Message   string         `json:"message,omitempty"`
	Error     string         `json:"error,omitempty"`
	Result    *TaskResult    `json:"result,omitempty"`
}

// NewEvent creates a new event with a timestamp.
func NewEvent(typ EventType, agentID, taskID string, turn int) *Event {
	return &Event{
		Type:      typ,
		AgentID:   agentID,
		TaskID:    taskID,
		Turn:      turn,
		Timestamp: time.Now(),
	}
}

// WithMessage sets the event message.
func (e *Event) WithMessage(msg string) *Event {
	e.Message = msg
	return e
}

// WithData sets arbitrary event data.
func (e *Event) WithData(data any) *Event {
	e.Data = data
	return e
}

// WithError sets the error string.
func (e *Event) WithError(err string) *Event {
	e.Error = err
	return e
}

// WithResult sets the task result.
func (e *Event) WithResult(r *TaskResult) *Event {
	e.Result = r
	return e
}

// ToolCallData holds data for a tool call event.
type ToolCallData struct {
	ToolName  string `json:"tool_name"`
	Arguments string `json:"arguments"`
	CallID    string `json:"call_id"`
}

// ToolResultData holds data for a tool result event.
type ToolResultData struct {
	ToolName string `json:"tool_name"`
	CallID   string `json:"call_id"`
	Result   any    `json:"result"`
	Error    string `json:"error,omitempty"`
}

// LLMCallData holds data about an LLM request.
type LLMCallData struct {
	Model        string `json:"model"`
	MessageCount int    `json:"message_count"`
	ToolCount    int    `json:"tool_count"`
}

// LLMResultData holds data about an LLM response.
type LLMResultData struct {
	FinishReason string `json:"finish_reason"`
	ToolCalls    int    `json:"tool_calls"`
	Tokens       int    `json:"tokens"`
}
