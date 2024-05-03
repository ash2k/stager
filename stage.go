package stager

import (
	"context"
)

type Stage interface {
	// Go starts f in a new goroutine attached to the Stage.
	// Stage context is passed to f as an argument. f should stop when context signals done.
	// If f returns a non-nil error, the stager starts performing shutdown.
	Go(f func(context.Context) error)
	// GoWhenDone starts f in a new goroutine attached to the Stage when the stage starts shutting down.
	// Stage shutdown waits for f to exit.
	GoWhenDone(f func() error)
}

type stage struct {
	ctx             context.Context
	cancelStage     context.CancelFunc
	cancelStagerRun context.CancelFunc
	whenDone        []func() error
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

func (s *stage) GoWhenDone(f func() error) {
	s.whenDone = append(s.whenDone, f)
}
