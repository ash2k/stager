package stager

import (
	"context"
)

type Stager interface {
	// NextStage adds a new stage to the Stager.
	NextStage() Stage
	// NextStageWithContext adds a new stage to the Stager. Provided ctxParent is used as the parent context for the
	// Stage's context.
	NextStageWithContext(ctxParent context.Context) Stage
	// Run blocks until ctx signals done or a function in a stage returns a non-nil error.
	// When it unblocks, it iterates Stages in reverse order. For each stage it cancels
	// it's context and waits for all started goroutines of that stage to finish.
	// Then it proceeds to the next stage.
	Run(ctx context.Context) error
}

func New() Stager {
	s := &stager{}
	s.runCtx, s.runCancel = context.WithCancel(context.Background())
	return s
}

type stager struct {
	stages    []*stage
	runCtx    context.Context
	runCancel context.CancelFunc
}

func (sr *stager) NextStage() Stage {
	return sr.NextStageWithContext(context.Background())
}

func (sr *stager) NextStageWithContext(ctxParent context.Context) Stage {
	ctx, cancel := context.WithCancel(ctxParent)
	st := &stage{
		ctx:             ctx,
		cancelStage:     cancel,
		cancelStagerRun: sr.runCancel,
		errChan:         make(chan error, 1),
	}
	sr.stages = append(sr.stages, st)
	return st
}

func (sr *stager) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
	case <-sr.runCtx.Done():
	}
	var firstErr error
	for i := len(sr.stages) - 1; i >= 0; i-- {
		st := sr.stages[i]
		st.cancelStage()
		for i := 0; i < st.n; i++ {
			err := <-st.errChan
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}
