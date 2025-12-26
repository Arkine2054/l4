package calendar

import (
	"fmt"
	"sync"
	"time"
)

type Service struct {
	storage    *Storage
	reminderCh chan ReminderJob

	muCancel  sync.Mutex
	cancelled map[int]bool
}

func NewService(st *Storage) *Service {
	return &Service{
		storage:    st,
		reminderCh: make(chan ReminderJob, 100),
		cancelled:  make(map[int]bool),
	}
}

func (s *Service) ReminderChannel() <-chan ReminderJob {
	return s.reminderCh
}

func (s *Service) Create(e Event) (Event, error) {
	s.storage.mu.Lock()
	e.ID = s.storage.nextID
	s.storage.nextID++
	s.storage.events[e.ID] = e
	s.storage.mu.Unlock()

	if e.RemindAt != nil {
		s.reminderCh <- ReminderJob{
			EventID:  e.ID,
			RemindAt: *e.RemindAt,
			Text:     e.Text,
		}
	}

	return e, nil
}

func (s *Service) Update(id int, e Event) error {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()

	if _, ok := s.storage.events[id]; !ok {
		return fmt.Errorf("event %d not found", id)
	}

	e.ID = id
	s.storage.events[id] = e
	return nil
}

func (s *Service) Delete(id int) error {
	s.storage.mu.Lock()
	if _, ok := s.storage.events[id]; !ok {
		s.storage.mu.Unlock()
		return nil
	}
	delete(s.storage.events, id)
	s.storage.mu.Unlock()

	s.muCancel.Lock()
	s.cancelled[id] = true
	s.muCancel.Unlock()

	return nil
}

func (s *Service) IsCancelled(id int) bool {
	s.muCancel.Lock()
	defer s.muCancel.Unlock()
	return s.cancelled[id]
}

func (s *Service) EventsForDay(userID int, date time.Time) []Event {
	s.storage.mu.RLock()
	defer s.storage.mu.RUnlock()

	var res []Event
	for _, e := range s.storage.events {
		if e.UserID == userID && sameDay(e.Date, date) {
			res = append(res, e)
		}
	}
	return res
}

func (s *Service) EventsForWeek(userID int, date time.Time) []Event {
	var res []Event
	start := date.AddDate(0, 0, -int(date.Weekday()))
	end := start.AddDate(0, 0, 7)
	s.storage.mu.RLock()
	defer s.storage.mu.RUnlock()
	for _, e := range s.storage.events {
		if e.UserID == userID && !e.Date.Before(start) && e.Date.Before(end) {
			res = append(res, e)
		}
	}
	return res
}

func (s *Service) EventsForMonth(userID int, date time.Time) []Event {
	var res []Event
	start := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 1, 0)
	s.storage.mu.RLock()
	defer s.storage.mu.RUnlock()
	for _, e := range s.storage.events {
		if e.UserID == userID && !e.Date.Before(start) && e.Date.Before(end) {
			res = append(res, e)
		}
	}
	return res
}

func (s *Service) CleanupOldEvents(olderThan time.Time) int {
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()

	removed := 0
	for id, e := range s.storage.events {
		if e.Date.Before(olderThan) {
			delete(s.storage.events, id)
			removed++
		}
	}
	return removed
}

func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
