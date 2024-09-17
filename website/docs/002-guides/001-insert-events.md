---
title: "Insert events"
sidebar_label: "Insert events"
sidebar_position: 1
---

The first step in working with KanthorQ is inserting events into the system, which will remain there for you to process later (until you delete them). To do this, you need to understand the event structure in KanthorQ. Once familiar with it, I'll show you how to insert events in a basic way, as well as how to do so transactionally.

## The event structure

An event in the KanthorQ system has the following structure. The most important properties you'll work with often are `Subject` and `Body`.

```go
type Event struct {
  Id        string   `json:"id" validate:"required"`
  Subject   string   `json:"subject" validate:"required,is_subject"`
  Body      []byte   `json:"body" validate:"required"`
  Metadata  Metadata `json:"metadata" validate:"required"`
  CreatedAt int64    `json:"created_at"`
}
```

### The `Subject`

The `Subject` property allows you to organize your events in a hierarchical structure. This concept is inspired by [NATS Subject-Based Messaging](https://docs.nats.io/nats-concepts/subjects). If you're familiar with RabbitMQ, it's similar to a [Routing Key](https://www.rabbitmq.com/tutorials/tutorial-five-go#topic-exchange).

For example, if you have an event for an order update, you might define the subject as `order.updated`. You can then write logic to handle this event.

As your business evolves, you might need to handle different versions of this event. You could either modify the existing logic or write new logic for the updated version.

- To support both old and new logic, you can define a new subject like `order.updated.v2`. You can filter all versions using the pattern `order.updated.>` so that both old and new events will be matched.

- Alternatively, to keep old and new versions separate, you could define a subject like `v2.order.updated`. Then `order.updated.>` will match the old version, and `v2.order.updated.>` will match the new one.

If your business expands to multiple regions, you could further specify subjects like `ap-southeast-1.order.created`, `ap-southeast-2.order.created`, and so on.

### The `Body`

The Body is an arbitrary byte array where you can store any kind of data. The most common use case is storing a JSON string, but you could also store binary data, such as images, or even encrypt the body before storing it.

For example:

```go
// pseudo code for demonstration

body, err := encrypt(...)
if err !=nil {
  log.Fatal(err)
}
events := []*entities.Event{
  entities.NewEvent("system.say_goodbye", body)),
}
pub.Send(ctx, events)

// Decrypt when receiving the event
sub.Receive(ctx, func(ctx context.Context, msg *entities.Message) error {
  data, err := decrypt(msg.Event.Body)
  if err !=nil {
    log.Fatal(err)
  }

  // Work with the decrypted data
})
```

### Other properties

- `Metadata`: An arbitrary map to store additional information about your event, like telemetry tracing.

- `Id`: A crucial property serving as the primary and partition key of the stream. It must be lexicographically sortable. Common options include [ULID](https://github.com/ulid/spec) and [KSUID](https://github.com/segmentio/ksuid).

:::info
KanthorQ is using `ULID` by default.
:::

## Inserting Events (Basic Way)

To simplify the process, KanthorQ provides helper methods for initializing both the publisher and the event.

To initialize a publisher, you need to define two options:

- `Connection`: The connection string for the PostgreSQL database.
- `StreamName`: The name of the stream where events will be stored. It serves same as [NATS JetStream Streams](https://docs.nats.io/nats-concepts/jetstream/streams) or [RabbitMQ Exchange](https://www.rabbitmq.com/tutorials/tutorial-three-go#exchanges).

```go
options := &publisher.Options{
  Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
  // Using default stream for demo
  StreamName: entities.DefaultStreamName,
}
// Initialize a publisher
pub, cleanup := kanthorq.Pub(ctx, options)
// Clean up after done
defer cleanup()
```

To initialize an event, define the `Subject` and `Body` using the `NewEvent` method:

```go
subject := "system.say_hello"
body := []byte("{\"msg\": \"Hello World!\"}")
event := entities.NewEvent(subject, body)
event.Metadata["version"] = "2"
event.Metadata["traceparent"] = "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
```

Bringing it all together:

```go
options := &publisher.Options{
  Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
  // Using default stream for demo
  StreamName: entities.DefaultStreamName,
}
// Initialize a publisher
pub, cleanup := kanthorq.Pub(ctx, options)
// Clean up after done
defer cleanup()

subject := "system.say_hello"
body := []byte("{\"msg\": \"Hello World!\"}")
event := entities.NewEvent(subject, body)
event.Metadata["version"] = "2"
event.Metadata["traceparent"] = "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"

// Publish events
events := []*entities.Event{event}
if err:= pub.Send(ctx, events); err != nil {
  // Handle error
}
```

## Inserting Events (Transactional Way)

One of the key features of KanthorQ is the ability to publish events transactionally, ensuring events are only published if the entire transaction is successful.

For example, when updating an order as the example we mentioned at the beginning, you can ensure that both the update and the event publishing are either fully successful or both fail.

```go
options := &publisher.Options{
  Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
  StreamName: entities.DefaultStreamName,
}
pub, cleanup := kanthorq.Pub(ctx, options)
defer cleanup()

subject := "order.updated"
body := []byte("{\"txn_id\": \"afe86f5d-66a0-49ca-8c18-fbea71dc2a98\"}")
event := entities.NewEvent(subject, body)
event.Metadata["traceparent"] = "00-80e1afed08e019fc1110464cfa66635c-7a085853722dc6d2-01"

events := []*entities.Event{event}

// ------------ THE DIFFERENT IS HERE ---------
// Start a new transaction
conn, err := pgx.Connect(ctx, cm.uri)
if err != nil {
  return nil, err
}
tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
if err != nil {
  return nil, err
}

// Publish events transactionally
if err:= pub.SendTx(ctx, events, tx); err != nil {
  // handle error
}

// do other stuff
// call tx.Rollback(ctx) to abort the transaction

// Commit the transaction
if err := tx.Commit(ctx); err != nil {
  // handle error
}
```

For a full example, see our documentation on [Transactional Publisher](https://github.com/kanthorlabs/kanthorq/blob/main/example/transactional-publisher/main.go)
