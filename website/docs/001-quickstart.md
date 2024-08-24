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

## Sending event with a publisher

```go
var DATABASE_URI = "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable"
var options = &kanthorq.PublisherOptions{
  Connection: DATABASE_URI,
  StreamName: kanthorq.DefaultStreamName
}

// init a context with 5 seconds timeout
ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
defer cancel()

// init a publisher
pub, cleanup := kanthorq.Pub(ctx, &publisher.Options{
  Connection: DATABASE_URI,
  StreamName: entities.DefaultStreamName,
})
defer cleanup()

// define an event
subject := "system.say_hello"
body := []byte("{\"msg\": \"Hello World!\"}")
event := entities.NewEvent(subject, body)

// sending it to the stream with name entities.DefaultStreamName
if err:= pub.Send(ctx, event); err != nil {
  // handle error
}
```

## Handling events with a subscriber

```go
var DATABASE_URI = "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable"

// listen for SIGINT and SIGTERM so if you press Ctrl-C you can stop the program
ctx, stop := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
defer stop()

var options = &kanthorq.SubscriberOptions{
  StreamName: kanthorq.DefaultStreamName,
  ConsumerName: kanthorq.DefaultConsumerName,
  ConsumerSubjectFilter: []string{"system.>"},
  ConsumerAttemptMax:        entities.DefaultConsumerAttemptMax,
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
err := kanthorq.Sub(ctx, options, func(ctx context.Context, event *kanthorq.Event) error {
  ts := time.UnixMilli(event.CreatedAt).Format(time.RFC3339)
  // print out recevied events
  fmt.Printf("RECEIVED: %s | %s | %s\n", event.Id, event.Subject, ts)
  return nil
})

// print out error if any
if err != nil && !errors.Is(err, context.Canceled) {
  log.Fatal(err)
}
```
