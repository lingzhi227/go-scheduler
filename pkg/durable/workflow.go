package durable

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lingzhi227/go-scheduler/pkg/agent"
)

// WorkflowEngine provides durable execution for agent tasks.
// In production, this wraps go-workflows; this implementation provides
// the interface and a simple in-process fallback.
type WorkflowEngine interface {
	// StartWorkflow begins a durable workflow for a task.
	StartWorkflow(ctx context.Context, task *agent.Task) (string, error)
	// GetStatus returns the status of a running workflow.
	GetStatus(ctx context.Context, workflowID string) (*WorkflowStatus, error)
	// Signal sends a signal to a running workflow.
	Signal(ctx context.Context, workflowID string, signal string, data any) error
}

// WorkflowStatus describes the current state of a workflow.
type WorkflowStatus struct {
	WorkflowID string `json:"workflow_id"`
	TaskID     string `json:"task_id"`
	Status     string `json:"status"` // running, completed, failed, cancelled
	Output     string `json:"output,omitempty"`
	Error      string `json:"error,omitempty"`
}

// InProcessEngine is a simple non-durable workflow engine for development.
type InProcessEngine struct {
	workflows map[string]*WorkflowStatus
	runFn     func(ctx context.Context, task *agent.Task) (string, error)
}

func NewInProcessEngine(runFn func(ctx context.Context, task *agent.Task) (string, error)) *InProcessEngine {
	return &InProcessEngine{
		workflows: make(map[string]*WorkflowStatus),
		runFn:     runFn,
	}
}

func (e *InProcessEngine) StartWorkflow(ctx context.Context, task *agent.Task) (string, error) {
	wfID := "wf-" + task.ID
	e.workflows[wfID] = &WorkflowStatus{
		WorkflowID: wfID,
		TaskID:     task.ID,
		Status:     "running",
	}

	go func() {
		output, err := e.runFn(ctx, task)
		if err != nil {
			e.workflows[wfID].Status = "failed"
			e.workflows[wfID].Error = err.Error()
			slog.Error("workflow failed", "workflow_id", wfID, "error", err)
		} else {
			e.workflows[wfID].Status = "completed"
			e.workflows[wfID].Output = output
		}
	}()

	return wfID, nil
}

func (e *InProcessEngine) GetStatus(_ context.Context, workflowID string) (*WorkflowStatus, error) {
	ws, ok := e.workflows[workflowID]
	if !ok {
		return nil, fmt.Errorf("workflow %q not found", workflowID)
	}
	return ws, nil
}

func (e *InProcessEngine) Signal(_ context.Context, workflowID string, signal string, data any) error {
	_, ok := e.workflows[workflowID]
	if !ok {
		return fmt.Errorf("workflow %q not found", workflowID)
	}
	slog.Info("workflow signal received", "workflow_id", workflowID, "signal", signal)
	return nil
}
