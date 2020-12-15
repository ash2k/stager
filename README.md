# Stager

Stager is a library to help write code where you are in control of start and shutdown of concurrent operations.
I.e. you know when goroutines start and stop and in which order.

An [example](example/main.go) is below:
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
```
Output:
```
2020/12/15 15:34:41 Start 3
2020/12/15 15:34:41 Start 1.2
2020/12/15 15:34:41 Start 1.1
2020/12/15 15:34:41 Start 2
2020/12/15 15:34:46 Stop 3
2020/12/15 15:34:46 Stop 2
2020/12/15 15:34:46 Stop 1.2
2020/12/15 15:34:46 Stop 1.1
2020/12/15 15:34:46 Exiting main
```

Note the following:

- Shutdown order is deterministic - 3, 2, and then 1.
- Shutdown order within a stage is not deterministic - 1.1 and 1.2 are not ordered.
