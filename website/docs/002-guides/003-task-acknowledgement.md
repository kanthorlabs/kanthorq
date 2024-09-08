---
title: "Task acknowledgement"
sidebar_label: "Task acknowledgement"
sidebar_position: 3
---

## Implicit acknowledgement

By default, all subscribers will use implicit acknowledgement. That means KanthorQ system will automatically acknowledge tasks if there are no returng error in your handler. Otherwise, it will mark the task is not acknowledged and it should be retried.

The handler interface is described like this

```go
type Handler func(ctx context.Context, msg *Message) error
```

And there is the `Message` struct which contains the event and task information.

```go
type Message struct {
  // the event you push into a stream
  Event *entities.Event
  // the task that was generated from the event
  // contains all necessary about your execution on the event
  Task  *entities.Task
}
```

## Explicit acknowledgement

Sometime, you want to acknowledge the task manually, commit the acknowledgement to the database along with your business logic. So either you do something successfully and commit it to be done or nothing will be commited at all.

The `Message` struct provide two methods `Ack` and `Nack` to acknowledge or no-acknowledge the message respectively. Both of them are safe to call multiple times and concurrently.

```go
func (msg *Message) Ack(ctx context.Context) error
// Nack requires one more parameter `reason`
// so that we can know why the message is nacked and you can retrie it later
func (msg *Message) Nack(ctx context.Context, reason error) error
```

:::danger
So what happen if `Ack` and `Nack` are not successful? You need to retry it manually by yourself to guarantee consistency across your application.
:::

And there is an demonstration of how to use `Ack` and `Nack` explicitly. See [Acknowledgement example](https://github.com/kanthorlabs/kanthorq/blob/main/example/acknowledgement/main.go) for complete code.

```go
func(ctx context.Context, msg *subscriber.Message) error {
  // someone say hello, accept it
  if msg.Event.Subject == "system.say_hello" {
    return msg.Ack(ctx)
  }
  // I will miss you don't want to say goodbye, not acknowledge it
  return msg.Nack(ctx, errors.New("not saying goodbye"))
}
```

## Transactional acknowledgement

KanthorQ takes advantages of PostgreSQL's ACID transactional model to provide many features that ensures data consistency across your application. Like previous article about [Insert events in a transactional way](./001-insert-events.md#insert-events-in-a-transactional-way), you can also acknowledge task in a transactional way.

```go
// start the transcation
tx, err := conn.Begin(ctx)
if err != nil {
  return err
}
// do some business logic with the transaction

kanthorq.Sub(ctx, options,	func(ctx context.Context, msg *subscriber.Message) error {
  if msg.Event.Subject == "system.say_hello" {
    // acknowledge it with the transaction
    return msg.AckTx(ctx, tx)
  }
  // or nack it, also with the transaction
  return msg.NackTx(ctx, errors.New("not saying goodbye"))
})
```
