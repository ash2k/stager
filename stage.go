package stager

import (
	"context"
)

type Stage interface {
	// Go starts f in a new goroutine attached to the Stage.
	// Stage context is passed to f as an argument. f should stop when context signals done.
	// If f returns a non-nil error, the stager starts performing shutdown.
	Go(f func(context.Context) error)
}

type stage struct {
	ctx             context.Context
	cancelStage     context.CancelFunc
	cancelStagerRun context.CancelFunc
	errChan         chan error
	n               int
}

func (s *stage) Go(f func(context.Context) error) {
	s.n++
	go func() {
		err := f(s.ctx)
		if err != nil {
			s.cancelStagerRun()
		}
		s.errChan <- err
	}()
}
