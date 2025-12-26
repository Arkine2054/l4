package logger

import (
	"fmt"
	"sync"
)

type Logger struct {
	mu sync.Mutex
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Log(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Println(msg)
}
