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
  kanthorq migrate up -s 'github://kanthorlabs/kanthorq/migration#main' -d 'postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable'
  ```

## Register producer

```go
var DATABASE_URI := "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable"

publisher, err := kanthorq.NewPublisher(DATABASE_URI)
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

topic := "testing.demo"
body := []byte("{\"ping\": true}")
if err:= publisher.Send(ctx, kanthorq.NewEvent(topic, body)); err != nil {
  // handle error
}
```

## Register subscriber

```go
var DATABASE_URI := "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable"

subscriber, err := kanthorq.NewSubscriber(DATABASE_URI)
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
