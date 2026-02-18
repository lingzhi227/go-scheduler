package memory

import (
	"context"
	"sync"

	"github.com/lingzhi227/go-scheduler/pkg/vllm"
)

// WorkingMemory holds the current session's conversation messages.
// Corresponds to MemGen's "working memory" tier.
type WorkingMemory struct {
	mu       sync.RWMutex
	messages []vllm.Message
	maxSize  int
}

func NewWorkingMemory(maxSize int) *WorkingMemory {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &WorkingMemory{
		messages: make([]vllm.Message, 0, maxSize),
		maxSize:  maxSize,
	}
}

// Append adds a message to the working memory, evicting the oldest if full.
func (w *WorkingMemory) Append(msg vllm.Message) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.messages) >= w.maxSize {
		// Evict oldest non-system message
		start := 0
		for i, m := range w.messages {
			if m.Role != "system" {
				start = i
				break
			}
		}
		if start < len(w.messages)-1 {
			w.messages = append(w.messages[:start], w.messages[start+1:]...)
		}
	}
	w.messages = append(w.messages, msg)
}

// Messages returns a copy of all messages.
func (w *WorkingMemory) Messages() []vllm.Message {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]vllm.Message, len(w.messages))
	copy(out, w.messages)
	return out
}

// Clear removes all messages.
func (w *WorkingMemory) Clear() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.messages = w.messages[:0]
}

// Len returns the number of messages.
func (w *WorkingMemory) Len() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.messages)
}

// Summary generates a short summary of working memory (last N messages).
func (w *WorkingMemory) Summary(_ context.Context) string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if len(w.messages) == 0 {
		return "empty"
	}
	last := w.messages[len(w.messages)-1]
	return last.Role + ": " + truncate(last.Content, 200)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
