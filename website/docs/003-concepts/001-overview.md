---
title: "Overview"
sidebar_label: "Overview"
sidebar_position: 1
---

Let's discover what the KanthorQ architecture is and how components communicate with each others.

## Architecture

The KanthorQ architecture consists of four key components:

- **Publisher**: Responsible for inserting events into the KanthorQ system. This can be done via a Command Line Tool or Golang code within your application.
- **Stream**: Receives events and persists them within the system, organized by subjects.
- **Consumer**: Stores tasks generated from events in the Stream. An event can create multiple tasks across different Consumers, but within a single Consumer, only one task can be tied to an event.
- **Subscriber**: Part of your application that retrieves tasks from the Consumer and executes business logic.

![KanthorQ Architecture](/003-concepts/001-overview/kanthorq-architecture.svg)

## The Publisher

The **Publisher** interacts with the KanthorQ Stream to insert events into the system. When initializing a Consumer, you must specify the associated Stream, so the system knows where to send the event.

As the **Publisher** is simply application code, it can be implemented using Go, a CLI, or even an HTTP request (coming soon).

## The Stream

The Stream is where events are stored, allowing you to retrieve them later for processing. It receives events from the **Publisher**, organizes them in a _time-series_ format, and retains them until explicitly removed.

A Stream in KanthorQ can store any type of event. For example, both internal and business-related events can be stored in a single stream named "default," but this may not be well-organized. We recommend defining specific Streams for different purposes. For instance:

The `order_update` Stream only contains events related to order statuses, such as `order.created`, `order.confirmed`, and `order.cancelled`.
The `parcel_update` Stream is for third-party logistics events, like `parcel.shipping`, `parcel.lost`, and `parcel.received`.

:::info

Since Streams are organized as time-series data, it’s best to query them using the timestamp column for optimal performance.

:::

:::tip

Events in Stream are sorted ascending by default because we use the [ULID](https://github.com/ulid/spec) as the primary key.

:::

Events in a Stream are categorized by subjects, which are dot-separated words. You can use this structure in various scenarios:

- `order.cancelled` and `order.created`: Different event types.
- `order.cancelled` and `v1.order.cancelled`: Events published by different codebases or versions.
- `order.cancelled` and `ap-southeast-1.order.cancelled`: Events categorized by region.
- `order.cancelled` and `tier-starter.order.cancelled`: Events distinguished by tier.

## The Consumer

When an event is published, it needs to be processed based on your business logic. The event processing could succeed or fail. If successful, the event is simply marked as completed, and you move to the next event. However, handling a failed event is more complex:

- Should the event be retried?
- When should the retry occur? 15 minutes? 30 minutes?
- How many retries should be attempted?
- What happens if retries are exhausted? Should the event be deleted?
- More and more question will be raised ...

A Consumer helps answer these questions by storing tasks generated from events, each with its own metadata. For example, if `order.cancelled` is triggered, you may want two separate actions: sending a notification email and handling refund processing. These actions can be managed by different Consumers, each with distinct retry logic. One Consumer could retry 10 times for refund processing, while another only retries 3 times for email notifications.

:::tip

Although each Consumer in the diagram handles a single subject, you can define a Consumer to handle multiple subjects, such as both `order.cancelled` and `order.failed` for sending customer apology emails.

:::

## The Subscriber

Finally, the Subscriber plays the most crucial role, executing the business logic for event tasks. The Subscriber pulls tasks from the Consumer, processes them, and updates the task metadata accordingly:

- If the task is successful, it is marked as `Completed`.
- If a retry is needed, the Subscriber sets a time for the next attempt.
- If retries are exhausted, the task is marked as `Discarded`.

Multiple Subscribers can handle tasks from a single Consumer, allowing parallel processing. For instance, in the diagram, _Email Subscriber I_ processes tasks with IDs 1 and 4, while _Email Subscriber II_ handles the task with ID 3.
