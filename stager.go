package stager

import "context"

type Stager interface {
	NextStage() Stage
	Shutdown()
}

func New() Stager {
	return &stager{}
}

type stager struct {
	stages []*stage
}

func (sr *stager) NextStage() Stage {
	ctx, cancel := context.WithCancel(context.Background())
	st := &stage{
		ctx:    ctx,
		cancel: cancel,
	}
	sr.stages = append(sr.stages, st)
	return st
}

func (sr *stager) Shutdown() {
	for i := len(sr.stages) - 1; i >= 0; i-- {
		st := sr.stages[i]
		st.cancel()
		st.group.Wait()
	}
}
