package worker

import "context"

type semaphore struct {
	ch chan struct{}
}

func newSemaphore(n int) *semaphore {
	return &semaphore{ch: make(chan struct{}, n)}
}

func (s *semaphore) acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *semaphore) release() {
	<-s.ch
}

// drain waits until all acquired slots have been released.
func (s *semaphore) drain() {
	n := cap(s.ch)
	for i := 0; i < n; i++ {
		s.ch <- struct{}{}
	}
	for i := 0; i < n; i++ {
		<-s.ch
	}
}
