---
title: "Task acknowledgement"
sidebar_label: "Task acknowledgement"
sidebar_position: 3
---

## Implicit Acknowledgement

By default, all subscribers in KanthorQ use implicit acknowledgement. This means that the system automatically acknowledges tasks if no error is returned from your handler. If an error occurs, the system will mark the task as not acknowledged, and it will be retried.

The handler interface is defined as follows:

```go
type Handler func(ctx context.Context, msg *Message) error
```

The `Message` struct, passed to the handler, contains information about the event and task:

```go
type Message struct {
  // the event you push into a stream
  Event *entities.Event
  // the task that was generated from the event
  // contains all necessary about your execution on the event
  Task  *entities.Task
}
```

## Explicit Acknowledgement

In some cases, you may want to acknowledge tasks manually—committing the acknowledgment to the database along with your business logic. This ensures that either everything is committed (including the acknowledgment) or nothing is committed, maintaining data consistency.

The `Message` struct provides two methods for manual acknowledgement:

```go
func (msg *Message) Ack(ctx context.Context) error
// Nack requires a reason parameter, so you can log why the task wasn't acknowledged
func (msg *Message) Nack(ctx context.Context, reason error) error
```

:::danger

**What happens if Ack or Nack fail?**

If you cannot acknowledge the message (which represents a task and a corresponding event) on time, the Availability Subscriber will pick it up and start processing it again if you have set that subscriber.

:::

Here's a demonstration of how to use Ack and Nack explicitly:

```go
func(ctx context.Context, msg *subscriber.Message) error {
  // Accept and acknowledge if the subject is "system.say_hello"
  if msg.Event.Subject == "system.say_hello" {
    if err := msg.Ack(ctx); err != nil {
      // Handle ack error
    }
  }
  // I will miss you don't want to say goodbye, not acknowledge it
  if msg.Event.Subject == "system.say_goodbye" {
    if err := msg.Nack(ctx, errors.New("not saying goodbye")); err != nil {
      // Handle nack error
    }
  }

  return nil
}
```

See [Acknowledgement example](https://github.com/kanthorlabs/kanthorq/blob/main/example/acknowledgement/main.go) for the complete code.

## Transactional Acknowledgement

KanthorQ leverages PostgreSQL’s ACID transactional model to ensure data consistency. This allows you to acknowledge tasks within a transaction, ensuring that both your business logic and the task acknowledgment are committed together.

```go
kanthorq.Sub(ctx, options, func(ctx context.Context, msg *subscriber.Message) error {
  // Begin a PostgreSQL transaction:
  tx, err := conn.Begin(ctx)
  if err != nil {
    return err
  }

  // Accept and acknowledge if the subject is "system.say_hello"
  if msg.Event.Subject == "system.say_hello" {
    if err := msg.AckTx(ctx, tx); err != nil {
      // Handle ack error
    }
  }
  // I will miss you don't want to say goodbye, not acknowledge it
  if msg.Event.Subject == "system.say_goodbye" {
    if err := msg.NackTx(ctx, errors.New("not saying goodbye"), tx); err != nil {
      // Handle nack error
    }
  }

  if err:=tx.Commit(ctx); err != nil {
    // Handle any errors that occur during commit
  }
})
```
