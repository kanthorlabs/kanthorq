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
docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=changemenowornever -d postgres:16
```

## Installation

To install KanthorQ, run the following in the directory of a Go project (where a go.mod file is present):

```
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
  kanthorq migrate up -s 'github://kanthorlabs/kanthorq/migration#main' -d 'postgres://postgres:changemenowornever@localhost:5432/postgres?sslmode=disable'
  ```

## Register producer

To start publishing events, you need to follow these steps

- Initialize a publisher publisher instance with given PostgreSQL connection string and a Stream name
- Start the instance so it will prepare everything: connect the database, register the stream for you
- Send an event that includes its subject and body
- Stop the instance if you don't need it anymore

Example:

```go
var DATABASE_URI = "postgres://postgres:changemenowornever@localhost:5432/postgres?sslmode=disable"
var options = &kanthorq.PublisherOptions{
  StreamName: kanthorq.DefaultStreamName
}

publisher, err := kanthorq.NewPublisher(DATABASE_URI, options)
if err != nil {
  log.Panicf("could not create new publisher because of %v", err)
}

ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
defer cancel()

if err:= publisher.Start(ctx); err != nil {
  log.Panicf("could not start the publisher because of %v", err)
}
defer func () {
  if err:= publisher.Stop(ctx); err != nil {
    log.Panicf("could not stop the publisher because of %v", err)
  }
}()

subject := "system.ping"
body := []byte("{\"alive\": true}")
if err:= publisher.Send(ctx, kanthorq.NewEvent(subject, body)); err != nil {
  // handle error
}
```

## Register subscriber

To subscribe to events in KanthorQ system, you need to

- Intialize a subscriber with PostgreSQL connection string and Consumer properties. The Consumer requires unique name and the listening subject (what you use to publish event before)
- Start the instance so it will prepare everything: connect the database, register the stream and consumer for you
- Start receiving events from KanthorQ system with given handler
- Stop the instance if you don't need it anymore

```go
var DATABASE_URI := "postgres://postgres:changemenowornever@localhost:5432/postgres?sslmode=disable"
var options = &kanthorq.SubscriberOptions{
  StreamName: kanthorq.DefaultStreamName,
  ConsumerName: "internal",
  ConsumerSubject: "system.ping",
  ConsumerAttemptMax: kanthorq.DefaultConsumerAttemptMax
}

subscriber, err := kanthorq.NewSubscriber(DATABASE_URI, options)
if err != nil {
  log.Panicf("could not create new subscriber because of %v", err)
}

ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
defer cancel()

if err:= subscriber.Start(ctx); err != nil {
  log.Panicf("could not start the subscriber because of %v", err)
}
defer func () {
  if err:= subscriber.Stop(ctx); err != nil {
    log.Panicf("could not stop the subscriber because of %v", err)
  }
}()

go subscriber.Receive(ctx, func(ctx context.Context, event *kanthorq.Event) error {
  // handle event logic
})

// listen for the cancellation signal.
<-ctx
```

If you feel the initialization of the subscriber includes some options that you do not fully understand, don't worry. Go to the next step, and I will explain it to you.
