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

	s := st.NextStage()
	s.Go(func(ctx context.Context) error {
		log.Print("Start 1.1")
		defer log.Print("Stop 1.1")
		<-ctx.Done()
		return nil
	})
	s.Go(func(ctx context.Context) error {
		log.Print("Start 1.2")
		defer log.Print("Stop 1.2")
		<-ctx.Done()
		return nil
	})

	s = st.NextStage()
	s.Go(func(ctx context.Context) error {
		log.Print("Start 2")
		defer log.Print("Stop 2")
		<-ctx.Done()
		return nil
	})

	s = st.NextStage()
	s.Go(func(ctx context.Context) error {
		log.Print("Start 3")
		defer log.Print("Stop 3")
		<-ctx.Done()
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := st.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
