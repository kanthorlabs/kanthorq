---
title: "Overview"
sidebar_label: "Overview"
sidebar_position: 1
---

Let's discover what the KanthorQ architecture is and how the components communicate with each other.

## Architecture

The KanthorQ architecture, depicted in the diagram below, comprises four components:

- **Publisher**: Application code responsible for inserting events into the KanthorQ system. This can be a Command Line Toolor Golang code within your application.
- **Stream**: Receives events and persists them within the KanthorQ system, categorized by subjects.
- **Consumer**: Stores tasks generated from events in the Stream. One event can generate many tasks in different _Consumers_, but in the same _Consumer_, there must be only one task that belongs to an event
- **Subscriber**: Part of your application that pulls tasks to execute your business logic.

![KanthorQ Architecture](/002-concepts/001-overview/kanthorq-architecture.svg)

## The Publisher

The Publisher is a component that interacts with KanthorQ Stream component to let you insert events into KanthorQ system. So that you need to specify the stream when you initialize the consumer to let it know where it should put the event into.

Because the Publisher is just an application code, so it can be your GO code, CLI or HTTP request (comming soon).

## The Stream

You need to store your events somewhere so you can retrieve it later to do your work. The Stream is that place! It will receive events from the Publisher, organize it as a **time-series** then store it until you explicitly remove it.

A Stream in KanthorQ system can store any kind of events. For example you can store all your internal events and your business events in the same stream with name `default`. It's okay but not well organized. That why we recommend you should define what a Stream should do then only put relative events into it. In the diagram we already showed you an example of how to organize a Stream and its events

- The `order_update` Stream will only receive events that are related to order changes like `order.created`, `order.confirmed`, `order.cancelled` and so on
- The `parcel_update` Stream for Third-party Logistics events like `parcel.shipping`, `parcel.lost`, `parcel.received` and so on

:::info

When we said that a Stream is time-series data, that mean you should only perform scan query on a Stream based on the timestamp column to achieve best performance.

:::

An event in a Stream will be categorized by `subject` what is dot-separated words. We can use it in various usecases an handle them differently

- `order.cancelled` and `order.created` is normal usecase that events are belonged to different type.
- `order.cancelled` and `v1.order.cancelled` indicate that events are published by different codebases that has different versions.
- `order.cancelled` and `ap-southeast-.order.cancelled` indicate that events are belonged to different regions.
- `order.cancelled` and `tier-starter.order.cancelled` indicate that events are tier-based.

## The Consumer

Whenever you published an event, you expect it should be handle later with your business logic. The action of handling that event can be successful or failed. If it is successful, it's easy because you may only want to mark that event as completed and move to next event. But failed processing invloves more complicated flow to handle it:

- If an event is failed to process, should we retry it?
- Retry? When? 15mins, 30mins or arbitrary value?
- If yes, how many time do we want to retry it?
- After configurable retries, the event is not successful yet, what next? Delete it?
- `Place here a thousand question about handling failed event...`

But it's not enough to introduce new component in our system because we can simply embedded metadata in the event itself to tell use what to do with failed event. The main question that leads to the born of The Consumer is

- Given an event, can I do more than one action independently for that event. For instance, if `order.cancelled` is fired, I want to not only sending our customer an nofication email about the cancellation but also do the refund cleanup works

Then we will have 2 separated Consumers that stores same event references but with different metadata about how events are handling. The same event with id `1` in refund handler can be retry 10 times but the one in email sending handler should only retry 3 times

:::tip

Although the diagram showed that each consumer contains only one subject, you can have as many subject as you want to have in a consumer. For example you can define a consumer that contains both `order.cancelled` and `order.failed` to send excuse emails to your customer.

:::

:::note

Friendly reminder that a Consumer stores many Tasks that are referenced to Events in a Stream

:::

## The Subscriber

Yoh we reach the final component that is the most important component because it contains your business logic to handle event tasks. The Subscriber need to know what Consumer it will subscribe to pull tasks from it, execute your business logic, then put back the metadata about that task

- If the task was executed successfully, mark the task as `Completed`
- If the task need retrying, the Subscriber need putting back what time it should be retried
- If the task retried many times and it finally reached the maximum retry times, the Subscriber should mark it as `Discarded`
- `Place here a thousand case about handling failed event...`

You can have many subscribers for a consumer as long as your system can handle system throughput. So that you can achieve Parallelism processing. In the diagram, we showed you that the _Email Subscriber I_ executed 2 tasks with id `1` and `3` and the _Email Subscriber II_ executed only one task with id `2`
