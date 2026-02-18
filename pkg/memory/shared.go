package memory

import (
	"context"
	"fmt"
	"sync"
)

// Permission defines access levels for shared spaces.
type Permission int

const (
	PermReader Permission = iota
	PermWriter
	PermAdmin
)

// SharedSpace allows agents to share data in named spaces with ACL.
type SharedSpace struct {
	mu      sync.RWMutex
	spaces  map[string]*spaceData
	backend Backend
}

type spaceData struct {
	members map[string]Permission // agentID -> permission
}

func NewSharedSpace(backend Backend) *SharedSpace {
	return &SharedSpace{
		spaces:  make(map[string]*spaceData),
		backend: backend,
	}
}

// CreateSpace creates a new shared space. The creator becomes admin.
func (s *SharedSpace) CreateSpace(space, creatorID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.spaces[space]; exists {
		return fmt.Errorf("space %q already exists", space)
	}
	s.spaces[space] = &spaceData{
		members: map[string]Permission{creatorID: PermAdmin},
	}
	return nil
}

// Join adds an agent to a space with the given permission.
func (s *SharedSpace) Join(space, agentID string, perm Permission) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sd, ok := s.spaces[space]
	if !ok {
		return fmt.Errorf("space %q not found", space)
	}
	sd.members[agentID] = perm
	return nil
}

// Leave removes an agent from a space.
func (s *SharedSpace) Leave(space, agentID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sd, ok := s.spaces[space]; ok {
		delete(sd.members, agentID)
	}
}

// Write writes a value to a shared space (requires writer or admin permission).
func (s *SharedSpace) Write(ctx context.Context, space, agentID, key string, value any) error {
	if err := s.checkPermission(space, agentID, PermWriter); err != nil {
		return err
	}
	return s.backend.Set(ScopeGlobal, space, key, value)
}

// Read reads a value from a shared space (any member can read).
func (s *SharedSpace) Read(ctx context.Context, space, agentID, key string) (any, bool, error) {
	if err := s.checkPermission(space, agentID, PermReader); err != nil {
		return nil, false, err
	}
	return s.backend.Get(ScopeGlobal, space, key)
}

// ListKeys lists all keys in a shared space.
func (s *SharedSpace) ListKeys(ctx context.Context, space, agentID string) ([]string, error) {
	if err := s.checkPermission(space, agentID, PermReader); err != nil {
		return nil, err
	}
	return s.backend.List(ScopeGlobal, space)
}

// ListSpaces returns all space names.
func (s *SharedSpace) ListSpaces() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.spaces))
	for name := range s.spaces {
		out = append(out, name)
	}
	return out
}

func (s *SharedSpace) checkPermission(space, agentID string, required Permission) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sd, ok := s.spaces[space]
	if !ok {
		return fmt.Errorf("space %q not found", space)
	}
	perm, ok := sd.members[agentID]
	if !ok {
		return fmt.Errorf("agent %q is not a member of space %q", agentID, space)
	}
	if perm < required {
		return fmt.Errorf("agent %q lacks permission in space %q (has %d, need %d)", agentID, space, perm, required)
	}
	return nil
}
