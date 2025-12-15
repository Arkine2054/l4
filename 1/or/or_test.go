package or

import (
	"testing"
	"time"
)

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func TestOrClosesOnFirstSignal(t *testing.T) {
	start := time.Now()

	<-Or(
		sig(100*time.Millisecond),
		sig(200*time.Millisecond),
		sig(300*time.Millisecond),
	)

	elapsed := time.Since(start)

	if elapsed > 150*time.Millisecond {
		t.Fatalf("or-channel closed too late: %v", elapsed)
	}
}

func TestOrWithSingleChannel(t *testing.T) {
	ch := sig(50 * time.Millisecond)

	start := time.Now()
	<-Or(ch)

	if time.Since(start) < 50*time.Millisecond {
		t.Fatal("or-channel closed too early")
	}
}

func TestOrWithNoChannels(t *testing.T) {
	if Or() != nil {
		t.Fatal("expected nil for zero channels")
	}
}
