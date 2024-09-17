---
title: "Quickstart"
sidebar_label: "Quickstart"
sidebar_position: 1
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

Learn how to install KanthorQ packages for Go, run migrations to set up KanthorQ's database schema, and start working with KanthorQ's publisher and subscriber.

## Prerequisites

To get started with KanthorQ, you only need one external service: a PostgreSQL database. However, you can use other databases that support the PostgreSQL wire protocol, such as CockroachDB or Amazon Aurora (PostgreSQL-compatible edition).

If you don't have a PostgreSQL instance, you can start one locally using Docker:

```bash
docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=changemenow -d postgres:16
```

## Installation

To install KanthorQ, run the following command in a Go project directory (where a `go.mod` file is present):

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
  Replace the `-d` option with your database URI if you're using a different instance.
  :::

## Sending Events with a Publisher

```go
import (
	"context"

	"github.com/kanthorlabs/kanthorq"
	"github.com/kanthorlabs/kanthorq/entities"
	"github.com/kanthorlabs/kanthorq/publisher"
)

func main() {
	ctx := context.Background()

	// Initialize a publisher
	pub, cleanup := kanthorq.Pub(ctx, &publisher.Options{
		// Replace the connection string with your database URI
		Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
		// Using the default stream for demo purposes
		StreamName: entities.DefaultStreamName,
	})
	// Clean up the publisher after use
	defer cleanup()

	subject := "system.say_hello"
	body := []byte("{\"msg\": \"Hello World!\"}")

	// Define your first event
	event := entities.NewEvent(subject, body)

	events := []*entities.Event{
		event,
		// Another  event
		entities.NewEvent("system.say_hello", []byte("{\"msg\": \"I'm comming!\"}")),
		// And yet another event
		entities.NewEvent("system.say_goodbye", []byte("{\"msg\": \"See you!!\"}")),
	}

	if err := pub.Send(ctx, events); err != nil {
		// Handle error
	}
}
```

## Handling Events with a Subscriber

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
	// Listen for SIGTERM, so pressing Ctrl-C stops the program
	ctx, stop := signal.NotifyContext(context.TODO(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var options = &subscriber.Options{
		// Replace the connection string with your database URI
		Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
		// Use the default stream for demo purposes
		StreamName: entities.DefaultStreamName,
		// Use the default consumer for demo purposes
		ConsumerName: entities.DefaultConsumerName,
		// Receive only events matching the filter,
		// so that both system.say_hello and system.say_goodbye will be processed
		ConsumerSubjectIncludes: []string{"system.>"},
		// Retry the task if it fails this many times
		ConsumerAttemptMax: entities.DefaultConsumerAttemptMax,
		// Reprocess stuck tasks after this duration
		ConsumerVisibilityTimeout: entities.DefaultConsumerVisibilityTimeout,
		Puller: puller.PullerIn{
			// Number of events to pull in one batch
			Size: 100,
			// Wait time before completing the batch if Size isn’t reached
			WaitingTime: 1000,
		},
	}

	// Handle events; this goroutine will block until Ctrl-C is pressed
	err := kanthorq.Sub(ctx, options, func(ctx context.Context, msg *subscriber.Message) error {
		ts := time.UnixMilli(msg.Event.CreatedAt).Format(time.RFC3339)
		// Print the received event
		fmt.Printf("RECEIVED: %s | %s | %s\n", msg.Event.Id, msg.Event.Subject, ts)
		return nil
	})

	// Print any errors, if applicable
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}

	fmt.Println("----- END OF EXAMPLE ------")
}
```

After running the example, you should see the following output:

```bash
2024/08/30 09:18:43 RECEIVED: event_01j6gh7t5v6j2q9ma0n78hw9fe | system.say_hello | 2024-08-30T09:18:42+07:00
2024/08/30 09:18:43 RECEIVED: event_01j6gh7t5s973x2sby12j9pwkc | system.say_hello | 2024-08-30T09:18:42+07:00
2024/08/30 09:18:45 RECEIVED: event_01j6gh7x3pvmk6demx3cq27j1q | system.say_goodbye | 2024-08-30T09:18:45+07:00
```

See the [Defaule Subscriber example](https://github.com/kanthorlabs/kanthorq/blob/main/example/default/main.go) for the complete code.
