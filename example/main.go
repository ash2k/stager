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
