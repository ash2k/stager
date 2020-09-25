package stager

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestShutdownOrder(t *testing.T) {
	var mx sync.Mutex
	var items []int
	st := New()

	s := st.NextStage()
	s.Go(func(ctx context.Context) error {
		<-ctx.Done()
		mx.Lock()
		defer mx.Unlock()
		items = append(items, 1)
		return nil
	})

	s = st.NextStage()
	s.Go(func(ctx context.Context) error {
		<-ctx.Done()
		mx.Lock()
		defer mx.Unlock()
		items = append(items, 2)
		return nil
	})

	s = st.NextStage()
	s.Go(func(ctx context.Context) error {
		<-ctx.Done()
		mx.Lock()
		defer mx.Unlock()
		items = append(items, 3)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := st.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(items, []int{3, 2, 1}) {
		t.Errorf("unexpected result %v", items)
	}
}

func TestShutdownOnError(t *testing.T) {
	var mx sync.Mutex
	var items []int
	st := New()

	s := st.NextStage()
	s.Go(func(ctx context.Context) error {
		mx.Lock()
		defer mx.Unlock()
		items = append(items, 1)
		return errors.New("boom")
	})

	s = st.NextStage()
	s.Go(func(ctx context.Context) error {
		<-ctx.Done()
		mx.Lock()
		defer mx.Unlock()
		items = append(items, 2)
		return nil
	})

	s = st.NextStage()
	s.Go(func(ctx context.Context) error {
		<-ctx.Done()
		mx.Lock()
		defer mx.Unlock()
		items = append(items, 3)
		return nil
	})

	err := st.Run(context.Background())
	if err == nil {
		t.Fatal("Expecting error")
	}
	if err.Error() != "boom" {
		t.Fatal("Expecting boom error")
	}

	if !reflect.DeepEqual(items, []int{1, 3, 2}) {
		t.Errorf("unexpected result %v", items)
	}
}

func TestRunReturnsFirstError(t *testing.T) {
	st := New()

	s := st.NextStage()
	s.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return errors.New("boom1")
	})

	s = st.NextStage()
	s.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return errors.New("boom2")
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := st.Run(ctx)
	if err == nil {
		t.Fatal("Expecting error")
	}
	if err.Error() != "boom2" {
		t.Fatal("Expecting boom2 error")
	}
}

func TestEmptyStagerStops(t *testing.T) {
	t.Run("no stages", func(t *testing.T) {
		st := New()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := st.Run(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("empty stage", func(t *testing.T) {
		st := New()
		st.NextStage()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := st.Run(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("stage with stopped func", func(t *testing.T) {
		st := New()

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		s := st.NextStage()
		s.Go(func(ctx context.Context) error {
			return nil
		})

		err := st.Run(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})
}
