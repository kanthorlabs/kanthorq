---
title: "Insert events"
sidebar_label: "Insert events"
sidebar_position: 1
---

The first step at a journey of working with KanthorQ is inserting events into KanthorQ system. So that it will be there for you to handle it later forever (until you delete it). To insert events into KanthorQ system, you need to know the structure of the events in KanthorQ system firstly. Then I will show you how to insert events in a basic way and in transactional way.

## The event structure

The representation of an event in KanthorQ system is shown below. The most common properties you will work with a lots later are the `Subject` and the `Body`.

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

The `Subject` is a property that allows you organize your events in a hierarchical structure. I borrowed it from [NATS Subject-Based Messaging](https://docs.nats.io/nats-concepts/subjects). If you are familiar with RabbitMQ, you can think it's kind of [Routing Key](https://www.rabbitmq.com/tutorials/tutorial-five-go#topic-exchange)

Lets say you need to work with event of order that is updated, you can define a subject like this `order.updated`. Then you can write your own logic to handle it.

After a period of time, your business requirements may change and you need another logic to handle the event that is belong. Then you have to decide whether you should modify the old logic to support the new version of just write another completely logic to handle new version.

- If you decide to support both old and new logic, you can define a new subject like this `order.updated.v2` so you can filter all events that match both the old and new version with the filter `order.updated.>`.
- If you choose the other soluton, you can define a new subject like this `v2.order.updated` so that the filter `order.updated.>` will only match the old one, and `v2.order.updated.>` only match the new one. Then you can register an other subscriber to working on the filter `v2.order.updated.>`

Later, you got a bump, your business grows rapidly and you decide that you need to support multiple regions. Then you can define another subject like this `ap-southeast-1.order.created`, `ap-southeast-2.order.created`, `ap-southeast-1.v2.order.created` and `ap-southeast-2.v2.order.created`

### The `Body`

The `body` is an arbitrary byte array that you can use to store any data you want. Most commone usage are json string, but you can also use image binary or base64 as well. And other example usage is you encrypt the `body` in the event before it's stored in the database. Then you can decrypt it when you need it.

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

// other logic

sub.Receive(ctx, func(ctx context.Context, msg *entities.Message) error {
  data, err := decrypt(msg.Event.Body)
  if err !=nil {
    log.Fatal(err)
  }

  // working with your data here
})
```

### Other properties

The `Metadata` is an arbitrary map that you can use to store additional information about your event. For example, KanthorQ will use it to store the Telemetry Tracing information

The `Id` is the most important property of an event but I think you should not touch it. It's a primary key as well as parition key of the stream. We use it to scan through the stream to get events so it must be lexicographically sortable. Some candidates that can be use here is

    - Our choice is [ULID](https://github.com/ulid/spec)
    - [KSUID](https://github.com/segmentio/ksuid). Uses UNIX-time in seconds, if your insert rate is about 1000 events per second, you loose the order of your inserting events.
    - Auto-increment ID of Postgres. Not unique if you try to looking in different streams as well as not available until you inserted it successfully.

## Insert events in a basic way

To make your coding experience a lot easier, I have define some facade methods to help you initialize both the publisher and event

To intialize the publisher, you must define two options:

    - The `Connection` is the connection string of PostgreSQl database
    - The `StreamName` is the name of the stream you want to store events. Think about [NATS JetStream Streams](https://docs.nats.io/nats-concepts/jetstream/streams) or [RabbitMQ Exchange](https://www.rabbitmq.com/tutorials/tutorial-three-go#exchanges)

```go
options := &publisher.Options{
  // replace connection string with your database URI
  Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
  // using default stream for demo
  StreamName: entities.DefaultStreamName,
}
// Initialize a publisher
pub, cleanup := kanthorq.Pub(ctx, options)
// clean up the publisher after everything is done
defer cleanup()
```

Initialize an event is easier, you only need to define the `Subject` and `Body` if you use the facade method `NewEvent`

```go
subject := "system.say_hello"
body := []byte("{\"msg\": \"Hello World!\"}")
event := entities.NewEvent(subject, body)
// add some additional metadata
event.Metadata["version"] = "2"
event.Metadata["traceparent"] = "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
```

Bring them up together, we will have a pseudo code like this

```go
options := &publisher.Options{
  // replace connection string with your database URI
  Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
  // using default stream for demo
  StreamName: entities.DefaultStreamName,
}
// Initialize a publisher
pub, cleanup := kanthorq.Pub(ctx, options)
// clean up the publisher after everything is done
defer cleanup()

subject := "system.say_hello"
body := []byte("{\"msg\": \"Hello World!\"}")
event := entities.NewEvent(subject, body)
// add some additional metadata
event.Metadata["version"] = "2"
event.Metadata["traceparent"] = "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"

// publish events
events := []*entities.Event{event}
if err:= pub.Send(ctx, events); err != nil {
  // handle error
}
```

## Insert events in a transactional way

One of the coolest features we have at KanthorQ is you can publish events in a transactional way. That means you can garantee that your events are only published if and only if the whole transcation is success. Is it cool, right?

If you get back to the example before about the order events, you can realise that you can expect that both your updating and the event publishing are either success or failure. You will not be able to fall into a case that you publish an event success but your updating is failed or vice versa.

```go
options := &publisher.Options{
  // replace connection string with your database URI
  Connection: "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable",
  // using default stream for demo
  StreamName: entities.DefaultStreamName,
}
// Initialize a publisher
pub, cleanup := kanthorq.Pub(ctx, options)
// clean up the publisher after everything is done
defer cleanup()

subject := "system.say_hello"
body := []byte("{\"msg\": \"Hello World!\"}")
event := entities.NewEvent(subject, body)
// add some additional metadata
event.Metadata["version"] = "2"
event.Metadata["traceparent"] = "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"

// publish events
events := []*entities.Event{event}
// ------------ THE DIFFERENT IS HERE ---------
// start a different connection
conn, err := pgx.Connect(ctx, cm.uri)
if err != nil {
  return nil, err
}
tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
if err != nil {
  return nil, err
}

if err:= pub.SendTx(ctx, events, tx); err != nil {
  // handle error
}
```

The full example can be found at our example of [Transactional Publisher](https://github.com/kanthor/kanthorq/blob/main/example/transactional-publisher/main.go)
