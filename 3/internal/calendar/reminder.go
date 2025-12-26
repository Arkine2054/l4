package calendar

import "time"

type ReminderJob struct {
	EventID  int
	RemindAt time.Time
	Text     string
}

type ReminderChecker interface {
	IsCancelled(id int) bool
}
