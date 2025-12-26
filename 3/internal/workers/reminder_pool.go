package workers

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/arkine/l4/3/internal/calendar"
)

func StartReminderPool(
	ctx context.Context,
	workers int,
	jobs <-chan calendar.ReminderJob,
	checker calendar.ReminderChecker,
	out func(string),
) {
	for i := 0; i < workers; i++ {
		go ReminderWorker(ctx, i, jobs, checker, out)
	}
}

func ReminderWorker(
	ctx context.Context,
	id int,
	jobs <-chan calendar.ReminderJob,
	checker calendar.ReminderChecker,
	out func(string),
) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-jobs:
			wait := time.Until(job.RemindAt)
			if wait > 0 {
				timer := time.NewTimer(wait)
				select {
				case <-ctx.Done():
					timer.Stop()
					return
				case <-timer.C:
				}
			}

			if checker != nil && checker.IsCancelled(job.EventID) {
				continue
			}

			if out != nil {
				out(job.Text)
			} else {
				fmt.Printf("[worker %d]  %s\n", id, job.Text)
			}
		}
	}
}
