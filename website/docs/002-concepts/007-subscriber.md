---
title: "Subscriber"
sidebar_label: "Subscriber"
sidebar_position: 7
---

The Subscriber is the most complicated component in KanthorQ system, but that complexity serves only one purpose: get your a task to work on it then try to get your task moves to **Final State**. If something went wrong with your task, you can ask for retry both manually or automatically from the Subscriber.

## Workflows

The Subscriber workflows will contains two parts: the pulling workflow that help you get tasks for your works and the updating workflow that help you update your task state after you have done with it

### The Pulling Workflow

```mermaid
---
title: Pulling flow
---
sequenceDiagram
  Subscriber ->> +Consumer Registry: name: send_cancellation_email

  rect rgb(191, 223, 255)
  note right of Subscriber: Transaction Box

  Consumer Registry ->> -kanthorq_stream_order_update: topic: order.cancelled, cursor: evt_01J36ZJACKR5FXDWVKASC4BNCN, limit: 100

  kanthorq_stream_order_update ->> kanthorq_stream_order_update: scanning
  kanthorq_stream_order_update ->> +send_cancellation_email: events
  send_cancellation_email -> -send_cancellation_email: convert events to tasks

  send_cancellation_email -->> +Subscriber: 100 tasks
  Subscriber ->> -Consumer Registry: next_cursor: evt_01J3702FVA6EJ7QB7CNRMCP93B

  end
```

Not like Publish only works with one component - the Stream, the Subscriber needs to interact with two components: the Stream and the Consumer. It will work with the Stream to help convert events from a stream to a task in a consumer, then it pulls those tasks for you. The _Transaction Box_ indicates that all actions will be run in a transcation, so that we can guarantee pulling a task exactly once.

1. We will start with a request to ask for 100 tasks.
2. We need to work with the Consumer Registry to get a stream name, a topic and a cursor of previous scanning in the Stream.
3. Put all parameters together we will scan the Stream to look for matching events with given topic.
4. After find events, we start converting those events to tasks by insert them into our Consumer then return those tasks back to our Subscriber.
5. Because a task is belong to only one event, so we also know what is the next cursor is (the latest task contains the latest matching event), so will update that cursor back to our Consumer Registry

:::info

By saying **scanning**, we mean we will query events from a stream from the lower bound that is specify by the **cursor** until we get enough rows (100 events). The simplify query will look like

```sql
SELECT * FROM kanthorq_stream_order_update WHERE id > 'evt_01J36ZJACKR5FXDWVKASC4BNCN' LIMIT 100
```

:::
