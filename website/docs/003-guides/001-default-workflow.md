---
title: "Default Workflow"
sidebar_label: "Default Workflow"
sidebar_position: 1
---

The default workflow of KanthorQ is simple, you insert your tasks with different subjects, then use consumer to filter which tasks you want to execute base on your subject filter. If everything works well, tasks will be marked as `Completed`, otherwise it will be `Retryable`.

```mermaid
---
title: Default Workflow
---
stateDiagram-v2
    direction LR
    [*] --> Available
    Available --> Running
    Running --> Completed
    Completed --> [*]
    Running --> Retryable
    Retryable --> [*]
```

:::warning
In this guideline, I will omit all error handling logic to make the code more readable. You should handle errors in your codebase to make sure it works well.
:::

Lets's start with the example codebase. You can find it in the [example/default](https://github.com/kanthorlabs/kanthorq/tree/main/example/default) folder.

```go
// import necessary packages
// import "github.com/kanthorq/kanthorq"

var DATABASE_URI = "postgres://postgres:changemenow@localhost:5432/postgres?sslmode=disable"

// Initialize a publisher
publisher, _ := kanthorq.NewPublisher(DATABASE_URI, &kanthorq.PublisherOptions{
  StreamName: kanthorq.DefaultStreamName,
})
publisher.Start(context.Background())
defer publisher.Stop(context.Background())

// publish some events, it will be stored inside our stream with name kanthorq.DefaultStreamName
publisher.Send(context.Background(), kanthorq.NewEvent("system.say_hello", []byte("{\"msg\": \"Hello World!\"}")))

// Initialize a subscriber that will process events that has subject that match with the filter "system.>"
// so both system.say_hello and system.say_goodbye will be processed
subscriber, _ := kanthorq.NewSubscriber(DATABASE_URI, &kanthorq.SubscriberOptions{
  StreamName: kanthorq.DefaultStreamName,
  ConsumerName: kanthorq.DefaultConsumerName,
  ConsumerSubjectFilter: "system.>",
  ConsumerAttemptMax: kanthorq.DefaultConsumerAttemptMax,
})
ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
defer cancel()
subscriber.Start(ctx)
defer subscriber.Stop(ctx)

go subscriber.Receive(ctx, func(ctx context.Context, event *kanthorq.Event) error {
  log.Print(string(event.Body))
  return nil
})

// wait for the cancellation signal
<-ctx.Done()
```

You may wonder what is subject filter (`system.>`) in the example codebase? Well, it is the subject filter that will be used to filter which events you want to process. And I inspired myself on the subject filter in the [Subject-Based Messaging](https://docs.nats.io/nats-concepts/subjects) of the project named [NATS.io](https://nats.io/).

I borrow the subject quote from the [NATS.io](https://nats.io/) project here

> What is a subject? At its simplest, a subject is just a string of characters that form a name the publisher and subscriber can use to find each other. More commonly subject hierarchies are used to scope messages into semantic namespaces.

You can find more information about how to use subject in KanthorQ at [Concepts - Subject](../002-concepts/008-subject.md)
