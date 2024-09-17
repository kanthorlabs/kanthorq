---
title: "Insert events"
sidebar_label: "Insert events"
sidebar_position: 1
---

The first step in working with the KanthorQ system is inserting events. Once an event is inserted, it will remain in the system indefinitely, ready to be processed at any time (unless you delete it). Before you can insert events, it’s important to understand the structure of events in the KanthorQ system. This section will cover the structure of events, how to insert them in a basic manner, and how to do so transactionally, ensuring that events are only added if the associated transaction succeeds.

## The event structure

An event in the KanthorQ system follows a predefined structure, as shown below. The most frequently used properties in the event structure are `Subject` and `Body`, which you will interact with often.

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

The Subject field is a crucial part of an event's structure. It allows you to organize your events in a hierarchical manner, similar to the concept of [NATS Subject-Based Messaging](https://docs.nats.io/nats-concepts/subjects). If you're familiar with RabbitMQ, you can think of it as being similar to a[Routing Key](https://www.rabbitmq.com/tutorials/tutorial-five-go#topic-exchange), which determines how messages (or in this case, events) are routed to different consumers.

For example, if you are working with events related to order updates, you can define a subject like `order.updated`. This allows you to easily organize all events that deal with order updates under a single subject. You can also define more granular subjects depending on your needs.

As your system evolves, you may need to introduce new logic to handle the events, perhaps as a result of changing business requirements. In such cases, you will need to decide whether to update the existing logic or create a new version of the event processing logic.

- If you choose to support both the old and the new logic simultaneously, you can define a new subject like `order.updated.v2`. This will allow you to filter both the old and new versions of the event using a single pattern, such as `order.updated.>`, which would match all versions of the `order.updated` subject.

- On the other hand, if you choose to keep the new logic separate from the old, you can define a subject like `v2.order.updated`. In this case, the filter `order.updated.>` would match only the old version, while `v2.order.updated.>` would match only the new one.

Furthermore, if your business grows and expands to multiple regions, you can organize your subjects by region. For example, you could define subjects like `ap-southeast-1.order.created` and `ap-southeast-2.order.created`, along with regional versions like `ap-southeast-1.v2.order.created `and `ap-southeast-2.v2.order.created`. This kind of flexibility allows you to organize and filter events in a way that suits the evolving structure of your business and its operational needs.

### The `Body`

The `Body` of an event is another important part of its structure. It is essentially an arbitrary byte array where you can store any kind of data you need. In most cases, developers use the `Body` field to store a JSON string, which represents structured data about the event. However, you are not limited to JSON; the `Body` field can also be used to store binary data, such as images, or even encoded or encrypted data, depending on your use case.

For example, you may choose to encrypt the data stored in the event body before saving it to the database. This approach ensures the security of your event data and can be useful in situations where sensitive information is involved. The data can then be decrypted when the event is consumed.

Here's an example of how you might encrypt the body of an event:

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

// Decrypt the data when receiving the event
sub.Receive(ctx, func(ctx context.Context, msg *entities.Message) error {
  data, err := decrypt(msg.Event.Body)
  if err !=nil {
    log.Fatal(err)
  }

  // Work with your decrypted data here
})
```

This flexibility in handling the `Body` of an event makes KanthorQ adaptable to a wide range of use cases, whether you need to work with simple JSON strings or more complex binary data formats.

### Other properties

In addition to the `Subject` and `Body` fields, events in KanthorQ have other properties that serve specific purposes:

- `Metadata`: This is an arbitrary map that can store additional information about the event. You can use this field to add any custom data related to the event. For instance, KanthorQ itself uses the `Metadata` field to store telemetry tracing information, which helps track the flow of events within distributed systems.

- `Id`: The `Id` field is a unique identifier for the event and plays a critical role in KanthorQ. It serves as both the primary key and the partition key within the event stream. This identifier must be lexicographically sortable, meaning the order in which events are inserted can be determined by their IDs. To ensure this, KanthorQ uses [ULID](https://github.com/ulid/spec) as the default method for generating event IDs, but other options are also available, such as [KSUID](https://github.com/segmentio/ksuid). However, ULID is preferred because it offers better guarantees for maintaining the correct order of events.

## Inserting Events (Basic Way)

To make event publishing easier, KanthorQ provides helper methods that simplify the process of initializing both the publisher and the event itself.

When initializing a publisher, you must define two key options:

- `Connection`: This is the connection string for the PostgreSQL database where events will be stored. You should replace this with the appropriate URI for your database.

- `StreamName`: This is the name of the stream in which you want to store your events. It’s akin to the concept of a stream in [NATS JetStream Streams](https://docs.nats.io/nats-concepts/jetstream/streams) or an exchange in [RabbitMQ Exchange](https://www.rabbitmq.com/tutorials/tutorial-three-go#exchanges).

```go
options := &publisher.Options{
  Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
  // Using default stream for demo
  StreamName: entities.DefaultStreamName, // Using the default stream for this example
}
// Initialize the publisher
pub, cleanup := kanthorq.Pub(ctx, options)
// Clean up the publisher after you're done
defer cleanup()
```

To initialize an event, you only need to define the `Subject` and `Body`. If you're using KanthorQ's helper methods, this process becomes even simpler. The `NewEvent` method can be used to create a new event, where the `Subject` describes the type of event, and the `Body` contains the event data (typically in JSON format).

```go
subject := "system.say_hello"
body := []byte("{\"msg\": \"Hello World!\"}")
event := entities.NewEvent(subject, body)

// Add some additional metadata to the event
event.Metadata["version"] = "2"
event.Metadata["traceparent"] = "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
```

Now, you can bring everything together by publishing the event:

```go
options := &publisher.Options{
  Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
  StreamName: entities.DefaultStreamName,
}
pub, cleanup := kanthorq.Pub(ctx, options)
defer cleanup()

subject := "system.say_hello"
body := []byte("{\"msg\": \"Hello World!\"}")
event := entities.NewEvent(subject, body)
event.Metadata["version"] = "2"
event.Metadata["traceparent"] = "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"

events := []*entities.Event{event}
if err:= pub.Send(ctx, events); err != nil {
  // Handle any errors that occur during event publishing
}
```

This is the basic process for inserting events into KanthorQ. It's simple but effective, allowing you to get up and running quickly.

## Inserting Events (Transactional Way)

One of the most powerful features of KanthorQ is its ability to handle event publishing in a transactional manner. This means you can ensure that events are only published if the entire transaction is successful. This feature is especially useful in scenarios where you need consistency between your business logic and event publishing.

For example, if you're updating an order in your system, you want to ensure that the order update and the event publication both either succeed or fail together. With transactional publishing, you can guarantee that no event is published unless the corresponding database transaction completes successfully.

Here’s how you can insert events transactionally:

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
  // Handle any errors that occur during transactional publishing
}

// Do whatever you need to do with the transaction
// call tx.Rollback(ctx) to abort the transaction

// Commit the transaction
if err := tx.Commit(ctx); err != nil {
  // Handle commit error
}
```

This code ensures that events are only published if the transaction completes successfully. If the transaction fails, the events will not be inserted into the stream. This feature provides a higher level of consistency and reliability in your event-driven system.

For more details and examples, refer to our full documentation on the [Transactional Publisher](https://github.com/kanthorlabs/kanthorq/blob/main/example/transactional-publisher/main.go)
