---
title: "Handle events and tasks"
sidebar_label: "Handle events and tasks"
sidebar_position: 2
---

Every event in KanthorQ generates at least one task. Depending on your requirements, you may need to handle potential failures, such as retrying tasks. This guide demonstrates two ways to manage tasks from your events using a subscriber. The first method is the most convenient, while the second provides more control over how to handle both success and failure scenarios.

## Using the Subscriber Facade

KanthorQ offers two facades: the **Publisher Facade** for quickly initializing your publisher to insert events into the KanthorQ system, and the **Subscriber Facade** to register and manage three types of subscribers:

- **Primary Subscriber** – Handles new incoming events and creates tasks to execute your business logic.
- **Retry Subscriber** – Manages events that require retries after a failure.
- **Visibility Subscriber** – Reprocesses tasks that exceed the visibility timeout and are considered "stuck" in the system.

Below is an example similar to the one in the Quickstart guide:

```go
import (
  "context"
  "errors"
  "fmt"
  "log"
  "os"
  "os/signal"
  "syscall"
  "time"

  "github.com/kanthorlabs/kanthorq"
  "github.com/kanthorlabs/kanthorq/entities"
  "github.com/kanthorlabs/kanthorq/puller"
  "github.com/kanthorlabs/kanthorq/subscriber"
)

func main() {
  // listen for SIGTERM so if you press Ctrl-C you can stop the program
  ctx, stop := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
  defer stop()

  var options = &subscriber.Options{
    Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
    StreamName: entities.DefaultStreamName,
    ConsumerName: entities.DefaultConsumerName,
    ConsumerSubjectIncludes: []string{"system.>"},
    ConsumerAttemptMax: entities.DefaultConsumerAttemptMax,
    ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
    Puller: puller.PullerIn{
      Size: 100,
      WaitingTime: 1000,
    },
  }

  // Handle events. This goroutine will block until Ctrl-C is pressed
   err := kanthorq.Sub(ctx, options, func(ctx context.Context, msg *subscriber.Message) error {
    ts := time.UnixMilli(msg.Event.CreatedAt).Format(time.RFC3339)
    // print out recevied event
    fmt.Printf("RECEIVED: %s | %s | %s\n", msg.Event.Id, msg.Event.Subject, ts)
    return nil
  })

  // print out error if any
  if err != nil && !errors.Is(err, context.Canceled) {
    log.Fatal(err)
  }
}
```

:::info
When using the **Subscriber Facade**, the handler you provide will be used for all subscribers under this facade. This means you’ll process tasks for new events, retrying events, and stuck events with the same logic in the handler.
:::

## Using the Subscriber Directly

In certain scenarios, you may want to handle tasks differently for new events and retries. In such cases, you can use the subscriber directly to achieve the desired behavior.

### Using the Primary Subscriber

This subscriber handles two main responsibilities:

- Scanning the stream for new events that match the filter set during consumer registration.
- Creating tasks corresponding to the events and returning them to you for execution.

Here’s an example of how to use the Primary Subscriber directly:

```go
import (
  "context"
  "errors"
  "fmt"
  "os"
  "os/signal"
  "syscall"
  "time"

  "github.com/kanthorlabs/kanthorq/entities"
  "github.com/kanthorlabs/kanthorq/pkg/xlogger"
  "github.com/kanthorlabs/kanthorq/puller"
  "github.com/kanthorlabs/kanthorq/subscriber"
  "go.uber.org/zap"
)

func main() {
  // listen for SIGTERM so if you press Ctrl-C you can stop the program
  ctx, stop := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
  defer stop()

  logger := xlogger.New()

  // options is same as the subscriber facade
  var options = &subscriber.Options{
    Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
    StreamName: entities.DefaultStreamName,
    ConsumerName: entities.DefaultConsumerName,
    ConsumerSubjectIncludes: []string{"system.>"},
    ConsumerAttemptMax: entities.DefaultConsumerAttemptMax,
    ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
    Puller: puller.PullerIn{
      Size: 100,
      WaitingTime: 1000,
    },
  }
  sub, err := subscriber.New(options, logger)
  if err != nil {
    panic(err)
  }

  var timeout = time.Second * 3

  // Starting a subscriber should be use with timeout
  startctx, cancel := context.WithTimeout(ctx, timeout)
  defer cancel()
  if err := sub.Start(startctx); err != nil {
    panic(err)
  }

  defer func() {
    // Graceful shutdown starting
    // don't reuse ctx here because it already done
    // you also need timeout here
    stopCtx, stopCancel := context.WithTimeout(context.Background(), timeout)
    defer stopCancel()
    if err := sub.Stop(stopCtx); err != nil {
      logger.Error("subscriber stop with error", zap.Error(err))
      return
    }
  }()

  // The main part, working on up comming events and tasks
  receiveCtx, receiveCancel := context.WithCancel(ctx)
  defer receiveCancel()

  // Start receiving events and tasks
  go func() {
    err := sub.Receive(receiveCtx, func(ctx context.Context, msg *subscriber.Message) error {
      ts := time.UnixMilli(msg.Event.CreatedAt).Format(time.RFC3339)
      // print out recevied event
      fmt.Printf("RECEIVED: %s | %s | %s\n", msg.Event.Id, msg.Event.Subject, ts)
      return nil
    })

    if err != nil && !errors.Is(err, context.Canceled) {
      logger.Error("subscriber receive with error", zap.Error(err))
    }

    // Subscriber is done, should cancel the context to trigger other workflows
    receiveCancel()
  }()

  <-receiveCtx.Done()
}
```

See the [Primary Subscriber example](https://github.com/kanthorlabs/kanthorq/blob/main/example/primary-subscriber/main.go) for the complete code.

### Using the Retry Subscriber

If an event task fails and is marked as retryable, the Retry Subscriber helps you retry it by transitioning the task from `Retryable` to `Running` and returning the task to you for execution. The code is similar to the **Primary Subscriber**, with only one line difference:

```git
- sub, err := subscriber.New(options, logger)
+ sub, err := subscriber.NewRetry(options, logger)
```

### Using the Availability Subscriber

For tasks that have exceeded the visibility timeout, the **Availability Subscriber** pulls them out for reprocessing. As with the **Retry Subscriber**, you only need to change one line:

```git
- sub, err := subscriber.New(options, logger)
+ sub, err := subscriber.NewAvailability(options, logger)
```
