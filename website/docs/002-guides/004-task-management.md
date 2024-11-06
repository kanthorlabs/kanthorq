---
title: "Task management"
sidebar_label: "Task management"
sidebar_position: 4
---

This guide introduces task management in KanthorQ, offering a hands-on look at interacting with KanthorQ’s core API. You'll get a clear view of how to use the API directly, allowing you to discover advanced ways to wokr with KanthorQ effectively.

## Cancellation

:::info
You can only mark a task as `Cancelled` if it's in `Pending`, `Available` or `Retryable` state.
:::
To cancel a task in KanthorQ, make sure you have the following ready:

- **Consumer**: Identify the consumer that the task belongs to.
- **PostgreSQL Connection**: Establish a connection using the `pgx` library.
- **Task**: Specify the task that you want to cancel.

These elements are essential for managing task cancellation directly through the KanthorQ system.

```go
// Assume `consumer` is already defined as a pointer to an entities.ConsumerRegistry struct

// Establish a connection to PostgreSQL
conn, err := pgx.Connect(ctx, DATABASE_URI)
if err != nil {
    // Handle connection error
    log.Fatalf("Failed to connect to database: %v", err)
}

// Define the cancellation request
cancellation, err := core.Do(ctx, conn, &core.TaskMarkCancelledReq{
    Consumer: consumer,
    Tasks:    []*entities.Task{task},
})

if err != nil {
    // Handle potential errors during task cancellation
    log.Fatalf("Failed to cancel task: %v", err)
}

// `cancellation` is a pointer to core.TaskMarkCancelledRes
// `Updated` contains the event IDs of tasks that have been successfully cancelled
fmt.Printf("Cancelled tasks with event IDs: %v\n", cancellation.Updated)

// `Noop` contains event IDs of tasks that couldn't be cancelled because they are:
// - Not in `Pending`, `Available`, or `Retryable` states
// - Not found in the registry with the given event ID
fmt.Printf("No operation occurred for tasks with event IDs: %v\n", cancellation.Noop)
```

See [Task Management example](https://github.com/kanthorlabs/kanthorq/blob/main/example/task-management/main.go) for the complete code.

## Resumption

:::info
You can only resume a task if it's in `Discarded` or `Cancelled` state.
:::
Once a task is cancelled, it can be resumed, allowing the system to process it again. The workflow for resuming a task closely resembles the cancellation process:

- Identify the consumer associated with the task.
- Establish a connection to PostgreSQL.
- Send a request to update the task's status to make it resumable.

Resuming a task involves modifying its state, thereby making it eligible for further processing.

```go
// Assume `consumer` is already defined as a pointer to an entities.ConsumerRegistry struct

// Establish a connection to PostgreSQL
conn, err := pgx.Connect(ctx, DATABASE_URI)
if err != nil {
    // Handle connection error
    log.Fatalf("Failed to connect to database: %v", err)
}

// Define the resumption request
resumption, err := core.Do(ctx, conn, &core.TaskResumeReq{
    Consumer: consumer,
    Tasks:    []*entities.Task{task},
})

if err != nil {
    // Handle potential errors during task resumption
    log.Fatalf("Failed to resume task: %v", err)
}

// `resumption` is a pointer to core.TaskResumeRes
// `Updated` contains the event IDs of tasks that have been successfully resumed
fmt.Printf("Resumed tasks with event IDs: %v\n", resumption.Updated)

// `Noop` contains event IDs of tasks that couldn't be resumed because they are:
// - Not in `Discarded` or `Cancelled` states
// - Not found in the registry with the given event ID
fmt.Printf("No operation occurred for tasks with event IDs: %v\n", resumption.Noop)
```
