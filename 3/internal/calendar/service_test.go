package calendar

import (
	"testing"
	"time"
)

func makeDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func TestCreateAndGetEvent(t *testing.T) {
	st := NewStorage()
	svc := NewService(st)
	date := makeDate(2025, 12, 26)
	ev, _ := svc.Create(Event{UserID: 1, Date: date, Text: "Test Event"})

	events := svc.EventsForDay(1, date)
	if len(events) != 1 || events[0].ID != ev.ID {
		t.Fatal("Event not found")
	}
}

func TestDeleteEventAndCancellation(t *testing.T) {
	st := NewStorage()
	svc := NewService(st)
	date := makeDate(2025, 12, 26)
	ev, _ := svc.Create(Event{UserID: 1, Date: date, Text: "ToDelete", RemindAt: &date})

	_ = svc.Delete(ev.ID)
	if !svc.IsCancelled(ev.ID) {
		t.Fatal("Event should be cancelled")
	}
}
