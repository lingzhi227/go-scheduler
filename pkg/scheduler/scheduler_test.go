package scheduler

import (
	"testing"
	"time"

	"github.com/lingzhi227/go-scheduler/pkg/agent"
)

func TestSchedulerPriority(t *testing.T) {
	// Use nil pool; we won't call Next (which needs a pool)
	s := NewScheduler(nil)

	s.Submit(&ScheduledTask{
		Task:        &agent.Task{ID: "low"},
		Priority:    1,
		SubmittedAt: time.Now(),
	})
	s.Submit(&ScheduledTask{
		Task:        &agent.Task{ID: "high"},
		Priority:    10,
		SubmittedAt: time.Now(),
	})
	s.Submit(&ScheduledTask{
		Task:        &agent.Task{ID: "med"},
		Priority:    5,
		SubmittedAt: time.Now(),
	})

	if s.QueueLen() != 3 {
		t.Fatalf("expected 3 tasks, got %d", s.QueueLen())
	}

	// Drain should return highest priority first
	tasks := s.Drain()
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}
	if tasks[0].Task.ID != "high" {
		t.Errorf("expected high priority first, got %s", tasks[0].Task.ID)
	}
	if tasks[1].Task.ID != "med" {
		t.Errorf("expected medium priority second, got %s", tasks[1].Task.ID)
	}
	if tasks[2].Task.ID != "low" {
		t.Errorf("expected low priority third, got %s", tasks[2].Task.ID)
	}

	if s.QueueLen() != 0 {
		t.Errorf("queue should be empty after drain, got %d", s.QueueLen())
	}
}

func TestSchedulerFIFOSamePriority(t *testing.T) {
	s := NewScheduler(nil)

	t1 := time.Now()
	t2 := t1.Add(time.Second)
	t3 := t2.Add(time.Second)

	s.Submit(&ScheduledTask{
		Task:        &agent.Task{ID: "first"},
		Priority:    5,
		SubmittedAt: t1,
	})
	s.Submit(&ScheduledTask{
		Task:        &agent.Task{ID: "third"},
		Priority:    5,
		SubmittedAt: t3,
	})
	s.Submit(&ScheduledTask{
		Task:        &agent.Task{ID: "second"},
		Priority:    5,
		SubmittedAt: t2,
	})

	tasks := s.Drain()
	if tasks[0].Task.ID != "first" {
		t.Errorf("expected first submitted first, got %s", tasks[0].Task.ID)
	}
}
