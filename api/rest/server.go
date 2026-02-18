package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	hollywoodactor "github.com/anthdm/hollywood/actor"
	"github.com/google/uuid"
	actorpkg "github.com/lingzhi227/go-scheduler/pkg/actor"
	"github.com/lingzhi227/go-scheduler/pkg/agent"
	"github.com/lingzhi227/go-scheduler/pkg/comm"
	"github.com/lingzhi227/go-scheduler/pkg/orchestrate"
)

// Server is the REST API server for go-scheduler.
type Server struct {
	engine     *hollywoodactor.Engine
	registry   *actorpkg.Registry
	supervisor *hollywoodactor.PID
	progress   *hollywoodactor.PID
	executor   *orchestrate.Executor

	mu     sync.RWMutex
	tasks  map[string]*TaskState
	mux    *http.ServeMux
	server *http.Server
}

// TaskState tracks a submitted task.
type TaskState struct {
	Task      *agent.Task    `json:"task"`
	Status    string         `json:"status"` // pending, running, completed, failed
	Result    *agent.TaskResult `json:"result,omitempty"`
	Error     string         `json:"error,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	DoneAt    *time.Time     `json:"done_at,omitempty"`
}

// NewServer creates a new REST API server.
func NewServer(
	engine *hollywoodactor.Engine,
	registry *actorpkg.Registry,
	supervisor *hollywoodactor.PID,
	progress *hollywoodactor.PID,
	executor *orchestrate.Executor,
) *Server {
	s := &Server{
		engine:     engine,
		registry:   registry,
		supervisor: supervisor,
		progress:   progress,
		executor:   executor,
		tasks:      make(map[string]*TaskState),
		mux:        http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("POST /tasks", s.handleSubmitTask)
	s.mux.HandleFunc("GET /tasks/{id}", s.handleGetTask)
	s.mux.HandleFunc("GET /tasks", s.handleListTasks)
	s.mux.HandleFunc("GET /agents", s.handleListAgents)
	s.mux.HandleFunc("GET /progress", s.handleGetProgress)
	s.mux.HandleFunc("POST /goals", s.handleSubmitGoal)
	s.mux.HandleFunc("GET /health", s.handleHealth)
}

// Start begins listening on the given address.
func (s *Server) Start(addr string) error {
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Minute, // long for LLM responses
	}
	slog.Info("REST API starting", "addr", addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// POST /tasks - Submit a task to a specific agent.
func (s *Server) handleSubmitTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Goal     string            `json:"goal"`
		AgentID  string            `json:"agent_id"`
		Context  string            `json:"context,omitempty"`
		MaxTurns int               `json:"max_turns,omitempty"`
		Input    map[string]any    `json:"input,omitempty"`
		Metadata map[string]string `json:"metadata,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}
	if req.Goal == "" {
		writeError(w, http.StatusBadRequest, "goal is required")
		return
	}

	task := &agent.Task{
		ID:       uuid.New().String(),
		Goal:     req.Goal,
		Context:  req.Context,
		MaxTurns: req.MaxTurns,
		Input:    req.Input,
		Metadata: req.Metadata,
	}

	// Find agent PID
	var agentPID *hollywoodactor.PID
	if req.AgentID != "" {
		pid, ok := s.registry.GetAgent(req.AgentID)
		if !ok {
			writeError(w, http.StatusNotFound, fmt.Sprintf("agent %q not found", req.AgentID))
			return
		}
		agentPID = pid
	} else {
		// Pick first available agent
		agents := s.registry.ListAgents()
		for _, pid := range agents {
			agentPID = pid
			break
		}
		if agentPID == nil {
			writeError(w, http.StatusServiceUnavailable, "no agents available")
			return
		}
	}

	// Track task state
	ts := &TaskState{
		Task:      task,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	s.mu.Lock()
	s.tasks[task.ID] = ts
	s.mu.Unlock()

	// Create a result collector actor
	collector := s.engine.SpawnFunc(func(ctx *hollywoodactor.Context) {
		switch msg := ctx.Message().(type) {
		case *actorpkg.TaskCompleted:
			s.mu.Lock()
			if ts, ok := s.tasks[msg.Result.TaskID]; ok {
				ts.Status = "completed"
				ts.Result = msg.Result
				now := time.Now()
				ts.DoneAt = &now
			}
			s.mu.Unlock()
		case *actorpkg.TaskFailed:
			s.mu.Lock()
			if ts, ok := s.tasks[msg.TaskID]; ok {
				ts.Status = "failed"
				ts.Error = msg.Error
				now := time.Now()
				ts.DoneAt = &now
			}
			s.mu.Unlock()
		}
	}, "collector")

	// Send task to agent
	ts.Status = "running"
	s.engine.Send(agentPID, &actorpkg.RunTask{
		Task:    task,
		ReplyTo: collector,
	})

	writeJSON(w, http.StatusAccepted, map[string]string{
		"task_id": task.ID,
		"status":  "running",
	})
}

// GET /tasks/{id}
func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.mu.RLock()
	ts, ok := s.tasks[id]
	s.mu.RUnlock()
	if !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, ts)
}

// GET /tasks
func (s *Server) handleListTasks(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	tasks := make([]*TaskState, 0, len(s.tasks))
	for _, ts := range s.tasks {
		tasks = append(tasks, ts)
	}
	s.mu.RUnlock()
	writeJSON(w, http.StatusOK, tasks)
}

// GET /agents
func (s *Server) handleListAgents(w http.ResponseWriter, r *http.Request) {
	agents := s.registry.ListAgents()
	type agentInfo struct {
		ID  string `json:"id"`
		PID string `json:"pid"`
	}
	var list []agentInfo
	for id, pid := range agents {
		list = append(list, agentInfo{ID: id, PID: pid.String()})
	}
	writeJSON(w, http.StatusOK, list)
}

// GET /progress
func (s *Server) handleGetProgress(w http.ResponseWriter, r *http.Request) {
	if s.progress == nil {
		writeJSON(w, http.StatusOK, map[string]any{"entries": map[string]any{}})
		return
	}
	resp := s.engine.Request(s.progress, &actorpkg.GetProgress{}, 5*time.Second)
	result, err := resp.Result()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get progress: "+err.Error())
		return
	}
	snapshot, ok := result.(*comm.ProgressSnapshot)
	if !ok {
		writeError(w, http.StatusInternalServerError, "unexpected progress response type")
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}

// POST /goals - Submit a high-level goal for decomposition and execution.
func (s *Server) handleSubmitGoal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Goal string `json:"goal"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Goal == "" {
		writeError(w, http.StatusBadRequest, "goal is required")
		return
	}

	// For now, create a simple single-agent task
	task := &agent.Task{
		ID:   uuid.New().String(),
		Goal: req.Goal,
	}

	ts := &TaskState{
		Task:      task,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	s.mu.Lock()
	s.tasks[task.ID] = ts
	s.mu.Unlock()

	writeJSON(w, http.StatusAccepted, map[string]string{
		"task_id": task.ID,
		"status":  "pending",
		"message": "goal submitted for decomposition",
	})
}

// GET /health
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

