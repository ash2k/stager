# Stager

Stager is a library to help write code where you are in control of start and shutdown of concurrent operations.
I.e. you know when goroutines start and stop and in which order.

An [example](example/main.go) is below. `pong` is always shutdown first, and `ping` last:
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

	ball := make(chan struct{})
	p1 := ping{
		ball: ball,
	}
	s := st.NextStage()
	s.Go(p1.run)

	p2 := pong{
		ball: ball,
	}
	s = st.NextStage()
	s.Go(p2.run)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := st.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

type ping struct {
	ball chan struct{}
}

func (p *ping) run(ctx context.Context) error {
	log.Print("Starting ping")
	defer log.Print("Shutting down ping")
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-p.ball:
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Second):
		}
		log.Print("ping")
		select {
		case <-ctx.Done():
			return nil
		case p.ball <- struct{}{}:
		}
	}
}

type pong struct {
	ball chan struct{}
}

func (p *pong) run(ctx context.Context) error {
	log.Print("Starting pong")
	defer time.Sleep(time.Second)
	defer log.Print("Shutting down pong - sleeping 1 second")
	for {
		select {
		case <-ctx.Done():
			return nil
		case p.ball <- struct{}{}:
		}
		select {
		case <-ctx.Done():
			return nil
		case <-p.ball:
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Second):
		}
		log.Print("pong")
	}
}
```
Output:
```bash
2017/06/20 13:34:33 Starting pong
2017/06/20 13:34:33 Starting ping
2017/06/20 13:34:34 ping
2017/06/20 13:34:35 pong
2017/06/20 13:34:36 ping
2017/06/20 13:34:37 pong
2017/06/20 13:34:38 Shutting down pong - sleeping 1 second
2017/06/20 13:34:38 ping
2017/06/20 13:34:39 Shutting down ping
2017/06/20 13:34:39 Exiting main
```
