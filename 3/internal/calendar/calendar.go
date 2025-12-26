package calendar

import (
	"sync"
	"time"
)

type Event struct {
	ID       int
	UserID   int
	Date     time.Time
	Text     string
	RemindAt *time.Time
}

type Storage struct {
	mu     sync.RWMutex
	events map[int]Event
	nextID int
}

func NewStorage() *Storage {
	return &Storage{
		events: make(map[int]Event),
		nextID: 1,
	}
}
