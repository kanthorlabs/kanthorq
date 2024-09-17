---
title: "Quickstart"
sidebar_label: "Quickstart"
sidebar_position: 1
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

To help you start working with KanthorQ, here's a guide on how to install the necessary packages, run database migrations, and begin publishing and subscribing to events in Go. This will walk you through setting up the core elements of KanthorQ and getting everything up and running.

## Prerequisites

Before diving into KanthorQ, you’ll need a PostgreSQL database. This can be a PostgreSQL instance running locally or in the cloud. Alternatively, you can use any database that supports the PostgreSQL wire protocol, such as CockroachDB or Amazon Aurora (PostgreSQL-compatible edition).

If you don’t have a PostgreSQL instance running, you can quickly start one locally using Docker:

```bash
docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=changemenow -d postgres:16
```

## Installation

To install KanthorQ, make sure you're in a Go project directory (one that contains a `go.mod` file). Then run the following command:

```bash
go get github.com/kanthorlabs/kanthorq
```

## Running migrations

KanthorQ relies on PostgreSQL to manage its events and tasks. To set up the necessary database schema, you’ll need to run some migrations. First, install the KanthorQ command-line tool:

```bash
go install github.com/kanthorlabs/kanthorq/cmd/kanthorq@latest
```

Next, run the migrations to set up KanthorQ’s database schema:

```bash
kanthorq migrate up -s 'github://kanthorlabs/kanthorq/migration#main' -d 'postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable'
```

Make sure to replace the -d option with the URI of your PostgreSQL instance if you're using a different database setup.

## Sending Events with a Publisher

Once the migration is complete, you’re ready to start sending events using the publisher. Here’s an example of how to publish events in Go:

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

In this example, you initialize a publisher that sends three events with different messages. The publisher handles event sending and interacts with the PostgreSQL database to persist those events.

## Handling Events with a Subscriber

Once you’ve sent some events, you’ll want to handle them using a subscriber. Here’s a basic example of how to subscribe to events:

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

This example shows a subscriber listening for events matching the subject filter `system.>`. The subscriber processes all events with subjects such as `system.say_hello` or `system.say_goodbye`.

After running the above example, you should see output similar to the following:

```bash
2024/08/30 09:18:43 RECEIVED: event_01j6gh7t5v6j2q9ma0n78hw9fe | system.say_hello | 2024-08-30T09:18:42+07:00
2024/08/30 09:18:43 RECEIVED: event_01j6gh7t5s973x2sby12j9pwkc | system.say_hello | 2024-08-30T09:18:42+07:00
2024/08/30 09:18:45 RECEIVED: event_01j6gh7x3pvmk6demx3cq27j1q | system.say_goodbye | 2024-08-30T09:18:45+07:00
```

See the [Defaule Subscriber example](https://github.com/kanthorlabs/kanthorq/blob/main/example/default/main.go) for the complete code.
