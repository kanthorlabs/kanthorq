---
title: "Quickstart"
sidebar_label: "Quickstart"
sidebar_position: 1
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

Learn how to install KanthorQ packages for Go, run migrations to get KanthorQ's database schema in place, and start working with KanthorQ publisher and subscriber.

## Prerequisites

To get started with KanthorQ you needs only one external service, a PostgreSQL database. But you can use others database that supports PostgreSQL wire protocol: CockroachDB or Amazon Aurora (PostgreSQL-compatible edition) for example.

If you don't have an instance of PostgreSQL, just start a new one in your local machine with Docker

```bash
docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=changemenow -d postgres:16
```

## Installation

To install KanthorQ, run the following in the directory of a Go project (where a go.mod file is present):

```bash
go get github.com/kanthorlabs/kanthorq
```

## Running migrations

KanthorQ system is replied on PosgreSQL database, and needs a small sets of tables to persist management and tasks data. You need to install the command line tool which executes migrations, and provides other features of KanthorQ system.

- Install the command line tool

  ```bash
  go install github.com/kanthorlabs/kanthorq/cmd/kanthorq@latest
  ```

- Run the migration up

  ```bash
  kanthorq migrate up -s 'github://kanthorlabs/kanthorq/migration#main' -d 'postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable'
  ```

  :::info
  Change the `-d` option with your database URI if you have a different database instance running rather than the default one
  :::

## Sending event with a publisher

```go
import "github.com/kanthorlabs/kanthorq"

func main() {
  // Initialize a publisher
  pub, cleanup := kanthorq.Pub(ctx, &publisher.Options{
    // replace DATABASE_URI with your database URI
    Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
    // using default stream for demo
    StreamName: entities.DefaultStreamName,
  })
  // clean up the publisher after everything is done
  defer cleanup()

  subject := "system.say_hello"
  body := []byte("{\"msg\": \"Hello World!\"}")
  // define your first event
  event := entities.NewEvent(subject, body)

  events:= []*entities.Event{
    event,
    // another event
    entities.NewEvent("system.say_hello", []byte("{\"msg\": \"I'm comming!\"}")),
    // and yet another event
    entities.NewEvent("system.say_goodbye", []byte("{\"msg\": \"See you!!\"}")),
  }

  if err:= pub.Send(ctx, events); err != nil {
    // handle error
  }
}
```

## Handling events with a subscriber

```go
import "github.com/kanthorlabs/kanthorq"

func main() {
  // listen for SIGINT and SIGTERM so if you press Ctrl-C you can stop the program
  ctx, stop := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
  defer stop()

  var options = &kanthorq.SubscriberOptions{
    // replace DATABASE_URI with your database URI
    Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
    // we use default stream for demo
    StreamName:                entities.DefaultStreamName,
    // we use default consumer for demo
    ConsumerName:              entities.DefaultConsumerName,
    // we will only receive events that match with the filter
    // so both system.say_hello and system.say_goodbye will be processed
    ConsumerSubjectFilter:     []string{"system.>"},
    // if task is failed, it will be retried it with this number of times
    ConsumerAttemptMax:        entities.DefaultConsumerAttemptMax,
    // if task is stuck, we will wait this amount of time to reprocess it
    ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
    Puller: &puller.PullerIn{
      // Size is how many events you want to pull at one batch
      Size:        100,
      // WaitingTime is how long you want to wait before pulling again
      // if you didn't get enough events in current batch
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

  fmt.Println("----- END OF EXAMPLE ------")
}
```

After running the example, you should see the following output:

```bash
2024/08/30 09:18:42 waiting for events...
2024/08/30 09:18:42 waiting for events...
2024/08/30 09:18:43 waiting for events...
2024/08/30 09:18:43 waiting for events...
2024/08/30 09:18:43 RECEIVED: event_01j6gh7t5v6j2q9ma0n78hw9fe | system.say_hello | 2024-08-30T09:18:42+07:00
2024/08/30 09:18:43 RECEIVED: event_01j6gh7t5s973x2sby12j9pwkc | system.say_hello | 2024-08-30T09:18:42+07:00
2024/08/30 09:18:43 waiting for events...
2024/08/30 09:18:43 waiting for events...
2024/08/30 09:18:44 waiting for events...
2024/08/30 09:18:44 waiting for events...
2024/08/30 09:18:44 waiting for events...
2024/08/30 09:18:44 waiting for events...
2024/08/30 09:18:44 waiting for events...
2024/08/30 09:18:45 waiting for events...
2024/08/30 09:18:45 waiting for events...
2024/08/30 09:18:45 waiting for events...
2024/08/30 09:18:45 waiting for events...
2024/08/30 09:18:45 RECEIVED: event_01j6gh7x3pvmk6demx3cq27j1q | system.say_goodbye | 2024-08-30T09:18:45+07:00
2024/08/30 09:18:46 waiting for events...
2024/08/30 09:18:46 waiting for events...
2024/08/30 09:18:46 waiting for events...
2024/08/30 09:18:46 waiting for events...
```

## Conclusion

By this tutorial, we showed you how quickly use Kanthorq in your project, the full and interactive example can be found in the [examples](https://github.com/kanthorlabs/kanthorq/blob/main/example/default/main.go) folder.
