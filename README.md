# Stager

Stager is a library to help write code where you are in control of start and shutdown of concurrent operations.
I.e. you know when goroutines start and stop and in which order.

A [runnable example](example/main.go):
```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/ash2k/stager"
)

func main() {
	defer log.Print("Exiting main")
	st := stager.New()
	defer st.Shutdown()

	ball := make(chan struct{})
	p1 := ping{
		ball: ball,
	}
	s := st.NextStage()
	s.StartWithContext(p1.run)

	p2 := pong{
		ball: ball,
	}
	s = st.NextStage()
	s.StartWithContext(p2.run)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	<-ctx.Done()
}

type ping struct {
	ball chan struct{}
}

func (p *ping) run(ctx context.Context) {
	log.Print("Starting ping")
	defer log.Print("Shutting down ping")
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.ball:
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
		}
		log.Print("ping")
		select {
		case <-ctx.Done():
			return
		case p.ball <- struct{}{}:
		}
	}
}

type pong struct {
	ball chan struct{}
}

func (p *pong) run(ctx context.Context) {
	log.Print("Starting pong")
	defer time.Sleep(time.Second)
	defer log.Print("Shutting down pong - sleeping 1 second")
	for {
		select {
		case <-ctx.Done():
			return
		case p.ball <- struct{}{}:
		}
		select {
		case <-ctx.Done():
			return
		case <-p.ball:
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
		}
		log.Print("pong")
	}
}
```