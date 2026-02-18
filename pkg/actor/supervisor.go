package actor

import (
	"log/slog"
	"time"

	hollywoodactor "github.com/anthdm/hollywood/actor"
)

// SupervisorStrategy determines how to handle child failures.
type SupervisorStrategy int

const (
	// OneForOne restarts only the failed child.
	OneForOne SupervisorStrategy = iota
	// AllForOne restarts all children when one fails.
	AllForOne
	// RestForOne restarts the failed child and all children started after it.
	RestForOne
)

// ChildSpec describes a child actor managed by the supervisor.
type ChildSpec struct {
	Kind     string
	Producer hollywoodactor.Producer
}

// SupervisorActor manages child actors with restart strategies.
type SupervisorActor struct {
	strategy    SupervisorStrategy
	maxRestarts int
	window      time.Duration
	specs       []ChildSpec
	children    map[string]*hollywoodactor.PID
	restarts    map[string][]time.Time
	registry    *Registry
}

// NewSupervisorProducer creates a Hollywood Producer for a SupervisorActor.
func NewSupervisorProducer(strategy SupervisorStrategy, specs []ChildSpec, registry *Registry) hollywoodactor.Producer {
	return func() hollywoodactor.Receiver {
		return &SupervisorActor{
			strategy:    strategy,
			maxRestarts: 5,
			window:      time.Minute,
			specs:       specs,
			children:    make(map[string]*hollywoodactor.PID),
			restarts:    make(map[string][]time.Time),
			registry:    registry,
		}
	}
}

func (s *SupervisorActor) Receive(ctx *hollywoodactor.Context) {
	switch msg := ctx.Message().(type) {
	case hollywoodactor.Started:
		s.startChildren(ctx)

	case hollywoodactor.Stopped:
		slog.Info("supervisor stopped")

	case *ChildFailed:
		s.handleChildFailure(ctx, msg)

	case *RestartChild:
		s.restartChild(ctx, msg.Kind)

	case *ListAgents:
		ctx.Respond(&AgentList{Agents: s.copyChildren()})
	}
}

func (s *SupervisorActor) startChildren(ctx *hollywoodactor.Context) {
	for _, spec := range s.specs {
		pid := ctx.SpawnChild(spec.Producer, spec.Kind)
		s.children[spec.Kind] = pid
		if s.registry != nil {
			s.registry.TrackAgent(spec.Kind, pid)
		}
		slog.Info("supervisor spawned child", "kind", spec.Kind)
	}
}

func (s *SupervisorActor) handleChildFailure(ctx *hollywoodactor.Context, msg *ChildFailed) {
	slog.Warn("child failed", "kind", msg.Kind, "error", msg.Error)

	if !s.canRestart(msg.Kind) {
		slog.Error("child exceeded max restarts, giving up", "kind", msg.Kind)
		return
	}

	switch s.strategy {
	case OneForOne:
		s.restartChild(ctx, msg.Kind)
	case AllForOne:
		s.restartAllChildren(ctx)
	case RestForOne:
		s.restartFromChild(ctx, msg.Kind)
	}
}

func (s *SupervisorActor) restartChild(ctx *hollywoodactor.Context, kind string) {
	// Find the spec
	for _, spec := range s.specs {
		if spec.Kind == kind {
			// Stop old child if still alive
			if old, ok := s.children[kind]; ok {
				ctx.Engine().Poison(old)
			}
			// Spawn new child
			pid := ctx.SpawnChild(spec.Producer, kind)
			s.children[kind] = pid
			s.recordRestart(kind)
			if s.registry != nil {
				s.registry.TrackAgent(kind, pid)
			}
			slog.Info("supervisor restarted child", "kind", kind)
			return
		}
	}
}

func (s *SupervisorActor) restartAllChildren(ctx *hollywoodactor.Context) {
	for _, spec := range s.specs {
		s.restartChild(ctx, spec.Kind)
	}
}

func (s *SupervisorActor) restartFromChild(ctx *hollywoodactor.Context, kind string) {
	found := false
	for _, spec := range s.specs {
		if spec.Kind == kind {
			found = true
		}
		if found {
			s.restartChild(ctx, spec.Kind)
		}
	}
}

func (s *SupervisorActor) canRestart(kind string) bool {
	now := time.Now()
	cutoff := now.Add(-s.window)
	restarts := s.restarts[kind]
	// Remove old restarts outside the window
	valid := restarts[:0]
	for _, t := range restarts {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	s.restarts[kind] = valid
	return len(valid) < s.maxRestarts
}

func (s *SupervisorActor) recordRestart(kind string) {
	s.restarts[kind] = append(s.restarts[kind], time.Now())
}

func (s *SupervisorActor) copyChildren() map[string]*hollywoodactor.PID {
	out := make(map[string]*hollywoodactor.PID, len(s.children))
	for k, v := range s.children {
		out[k] = v
	}
	return out
}
