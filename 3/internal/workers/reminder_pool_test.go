package workers

import (
	"context"
	"testing"
	"time"

	"gitlab.com/arkine/l4/3/internal/calendar"
)

type mockChecker struct {
	cancelled map[int]bool
}

func (m *mockChecker) IsCancelled(id int) bool {
	return m.cancelled[id]
}

func TestReminderWorkerExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan calendar.ReminderJob, 1)
	called := make(chan string, 1)
	out := func(s string) { called <- s }

	mock := &mockChecker{cancelled: map[int]bool{}}

	go ReminderWorker(ctx, 0, jobs, mock, out)

	now := time.Now()
	jobs <- calendar.ReminderJob{EventID: 1, RemindAt: now, Text: "test"}

	select {
	case res := <-called:
		if res != "test" {
			t.Fatalf("unexpected output: %s", res)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestReminderWorkerCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan calendar.ReminderJob, 1)
	called := make(chan string, 1)
	out := func(s string) { called <- s }

	mock := &mockChecker{cancelled: map[int]bool{1: true}}

	go ReminderWorker(ctx, 0, jobs, mock, out)

	now := time.Now()
	jobs <- calendar.ReminderJob{EventID: 1, RemindAt: now, Text: "should not fire"}

	select {
	case <-called:
		t.Fatal("reminder should have been cancelled")
	case <-time.After(100 * time.Millisecond):
		// OK
	}
}
