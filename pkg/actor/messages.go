package actor

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/lingzhi227/go-scheduler/pkg/agent"
)

// RunTask tells an AgentActor to start working on a task.
type RunTask struct {
	Task    *agent.Task
	ReplyTo *actor.PID // where to send TaskResult
}

// TaskCompleted is sent when an agent finishes a task.
type TaskCompleted struct {
	Result *agent.TaskResult
}

// TaskFailed is sent when an agent fails a task.
type TaskFailed struct {
	TaskID  string
	AgentID string
	Error   string
}

// ToolCallRequest asks a ToolExecutor to call a tool.
type ToolCallRequest struct {
	CallID    string
	ToolName  string
	Arguments []byte
	ReplyTo   *actor.PID
}

// ToolCallResponse is the result of a tool execution.
type ToolCallResponse struct {
	CallID   string
	ToolName string
	Result   any
	Error    string
}

// AgentProgress is sent to the ProgressBoard.
type AgentProgress struct {
	AgentID        string
	TaskID         string
	Status         string // running, done, error
	Turn           int
	Summary        string
	PartialResults []string
}

// GetProgress requests the current progress snapshot.
type GetProgress struct{}

// ProgressSnapshot is the response to GetProgress.
type ProgressSnapshot struct {
	Entries map[string]*AgentProgress
}

// ChildFailed notifies a supervisor that a child has failed.
type ChildFailed struct {
	ChildPID *actor.PID
	Kind     string
	Error    string
}

// RestartChild tells a supervisor to restart a specific child.
type RestartChild struct {
	Kind string
}

// ListAgents requests a list of all agent PIDs from the registry.
type ListAgents struct{}

// AgentList is the response to ListAgents.
type AgentList struct {
	Agents map[string]*actor.PID
}
