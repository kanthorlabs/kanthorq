---
title: "Working on events"
sidebar_label: "Working on events"
sidebar_position: 2
---

Every event will generated at least one task when you are working on it. Then, based on your expectation, you may need to work extra on that event to handle failure like retrying. In this article I will show you two ways to use subscriber to handle your tasks of your events. The first one is the most convivenient way. The second one gives your more control on how do you want to handle both happy case and failure case.

## Using subscriber facade

KanthorQ package provides you two facades, the first one is the publisher facade that you quickly intialize your publisher to help you insert your events into the KanthorQ system. Now let me introduce you to the second one - the subscriber facade that registers three types of subscribers:

- The Primary Subscriber that handles up comming events in the system, produces tasks to execute your business logic.
- The Retry Subscriber that handles events that need retrying.
- The Visibility Subscriber that handles events that are visible in the system, aka stuck for a long time in the system exceeded the visibility timeout.

The example code is similar to the one in the quickstart article.

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
    // replace connection string with your database URI
    Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
    // we use default stream for demo
    StreamName: entities.DefaultStreamName,
    // we use default consumer for demo
    ConsumerName: entities.DefaultConsumerName,
    // we will only receive events that match with the filter
    // so both system.say_hello and system.say_goodbye will be processed
    ConsumerSubjectIncludes: []string{"system.>"},
    // if task is failed, it will be retried it with this number of times
    ConsumerAttemptMax: entities.DefaultConsumerAttemptMax,
    // if task is stuck, we will wait this amount of time to reprocess it
    ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
    Puller: puller.PullerIn{
      // Size is how many events you want to pull at one batch
      Size: 100,
      // WaitingTime is how long you want to wait before finish current batch
      // because you don't get enough events defined in the Size attribute
      WaitingTime: 1000,
    },
  }

  // hanlding events, the gorouting will be block until you press Ctrl-C
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
Because you are using the subscriber facade, so that the handle you given will be use for all of subscribers under the facade. That mean you will process tasks of new events, retrying events and handling events that are stuck for a long time with same logic by the given handler.
:::

## Using subscriber directly

There are some usecases that you want to execute your task differently for the new events or the retry ones. In this case, you can use the subscriber directly to sastify your needs.

### Using the Primary Subscriber

This subscriber need working on two things

- Scan through your stream to look for new events that match with the filter you set when registering the consumer that will be used by the subscriber
- If there is new events, generate a task corresponding to that event, return both the task and the event to you so that you can execute it

The example code bellow shows how to use the Primary Subscriber directly:

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
    // replace connection string with your database URI
    Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
    // we use default stream for demo
    StreamName: entities.DefaultStreamName,
    // we use default consumer for demo
    ConsumerName: entities.DefaultConsumerName,
    // we will only receive events that match with the filter
    // so both system.say_hello and system.say_goodbye will be processed
    ConsumerSubjectIncludes: []string{"system.>"},
    // if task is failed, it will be retried it with this number of times
    ConsumerAttemptMax: entities.DefaultConsumerAttemptMax,
    // if task is stuck, we will wait this amount of time to reprocess it
    ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
    Puller: puller.PullerIn{
      // Size is how many events you want to pull at one batch
      Size: 100,
      // WaitingTime is how long you want to wait before finish current batch
      // because you don't get enough events defined in the Size attribute
      WaitingTime: 1000,
    },
  }
  sub, err := subscriber.New(options, logger)
  if err != nil {
    panic(err)
  }

  var timeout = time.Second * 3

  // starting a subscriber should be use with timeout
  startctx, cancel := context.WithTimeout(ctx, timeout)
  defer cancel()
  if err := sub.Start(startctx); err != nil {
    panic(err)
  }

  defer func() {
    // graceful shutdown starting
    // don't reuse ctx here because it already done
    // you also need timeout here
    stopCtx, stopCancel := context.WithTimeout(context.Background(), timeout)
    defer stopCancel()
    if err := sub.Stop(stopCtx); err != nil {
      logger.Error("subscriber stop with error", zap.Error(err))
      return
    }
  }()

  // the main part, working on up comming events and tasks
  receiveCtx, receiveCancel := context.WithCancel(ctx)
  defer receiveCancel()

  // start receiving events and tasks
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

    // subscriber is done, should cancel the context to trigger other workflows
    receiveCancel()
  }()

  <-receiveCtx.Done()
}
```

### Using the Retry Subscriber

If a task of event is failed and you mark it as retryable, this subscriber will help you retry it. It will

- Transition the state of task from `Retryable` to `Running`
- Return tasks that transitioned successfully to you to execute your business logic

The different between the Primary Subscriber and the Retry Subscriber in your point of view is small, you only need to change one line of code to get your Retry Subscriber to work.

```go
  // other codes is same as the primary subscriber

  // ------------ THE DIFFERENT IS HERE ---------
  // use subscriber.NewRetry instead of subscriber.New
  sub, err := subscriber.NewRetry(options, logger)
  if err != nil {
    panic(err)
  }

  // other codes is same as the primary subscriber
```

### Using the Availability Subscriber

If a task of event stays in the consumer for a long time and exceeds the visibility timeout, the Availability Subscriber will help you pull it out and work on it again.

- Set the new visibility time tasks
- Return tasks that has updated successfully to you to execute your business logic

The different between the Primary Subscriber and the Retry Subscriber in your point of view is small, you only need to change one line of code to get your Retry Subscriber to work.

```go
  // other codes is same as the primary subscriber

  // ------------ THE DIFFERENT IS HERE ---------
  // use subscriber.NewRetry instead of subscriber.New
  sub, err := subscriber.NewAvailability(options, logger)
  if err != nil {
    panic(err)
  }

  // other codes is same as the primary subscriber
```
