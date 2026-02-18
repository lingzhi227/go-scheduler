package scheduler

import (
	"container/heap"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/lingzhi227/go-scheduler/pkg/agent"
	"github.com/lingzhi227/go-scheduler/pkg/vllm"
)

// ScheduledTask wraps a task with scheduling metadata.
type ScheduledTask struct {
	Task            *agent.Task
	Priority        int    // higher = more urgent
	EstimatedTokens int    // estimated token budget
	Model           string // preferred model
	AgentID         string // assigned agent
	SubmittedAt     time.Time
	index           int // heap index
}

// Scheduler manages a priority queue of tasks and assigns them to vLLM servers.
type Scheduler struct {
	mu   sync.Mutex
	pool *vllm.Pool
	pq   priorityQueue
	cond *sync.Cond
}

func NewScheduler(pool *vllm.Pool) *Scheduler {
	s := &Scheduler{
		pool: pool,
		pq:   make(priorityQueue, 0),
	}
	s.cond = sync.NewCond(&s.mu)
	heap.Init(&s.pq)
	return s
}

// Submit adds a task to the scheduling queue.
func (s *Scheduler) Submit(task *ScheduledTask) error {
	if task.Task == nil {
		return fmt.Errorf("task is nil")
	}
	if task.SubmittedAt.IsZero() {
		task.SubmittedAt = time.Now()
	}
	s.mu.Lock()
	heap.Push(&s.pq, task)
	s.mu.Unlock()
	s.cond.Signal()
	return nil
}

// Next blocks until a task is available, then returns it with a selected vLLM client.
func (s *Scheduler) Next(ctx context.Context) (*ScheduledTask, *vllm.PoolClient, error) {
	s.mu.Lock()
	for s.pq.Len() == 0 {
		// Wait with context awareness
		done := make(chan struct{})
		go func() {
			s.cond.Wait()
			close(done)
		}()
		s.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-done:
		}
		s.mu.Lock()
	}

	task := heap.Pop(&s.pq).(*ScheduledTask)
	s.mu.Unlock()

	// Select a vLLM server based on task requirements
	req := &vllm.ChatCompletionRequest{
		Model: task.Model,
	}
	pc, err := s.pool.Get(ctx, req)
	if err != nil {
		// Re-queue on failure
		s.mu.Lock()
		heap.Push(&s.pq, task)
		s.mu.Unlock()
		return nil, nil, fmt.Errorf("no vllm server available: %w", err)
	}

	return task, pc, nil
}

// QueueLen returns the current queue length.
func (s *Scheduler) QueueLen() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pq.Len()
}

// Drain returns and removes all queued tasks.
func (s *Scheduler) Drain() []*ScheduledTask {
	s.mu.Lock()
	defer s.mu.Unlock()
	tasks := make([]*ScheduledTask, s.pq.Len())
	for i := 0; s.pq.Len() > 0; i++ {
		tasks[i] = heap.Pop(&s.pq).(*ScheduledTask)
	}
	return tasks
}

// priorityQueue implements heap.Interface for ScheduledTask.
type priorityQueue []*ScheduledTask

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// Higher priority first; on tie, earlier submission first
	if pq[i].Priority != pq[j].Priority {
		return pq[i].Priority > pq[j].Priority
	}
	return pq[i].SubmittedAt.Before(pq[j].SubmittedAt)
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*ScheduledTask)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[:n-1]
	return item
}
