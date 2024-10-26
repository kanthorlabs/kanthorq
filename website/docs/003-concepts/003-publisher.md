---
title: "Publisher"
sidebar_label: "Publisher"
sidebar_position: 3
---

The **Publisher** is responsible for inserting events into the KanthorQ system. Technically, it functions as a simple query that inserts events—nothing particularly special happens at this stage. However, when dealing with the insertion of multiple events in a short time, there are important performance considerations worth discussing.

## Basic usecase

In message brokers or queue systems, it's common to push a small number of events at a time. Here are a couple of typical scenarios:

- **Webhook Events**: These are processed one by one and pushed into the queue individually.
- **E-commerce Systems**: When processing an order, related events such as `order.updated` and `payment.initialized` are published. For each order, typically a few events are triggered

In these basic cases, the performance of the KanthorQ system is limited by PostgreSQL's transactions-per-second (TPS) capacity. Given this limitation, there's not much room for optimization in these scenarios.

## Batching usecase

Things become more complex when you need to insert many events simultaneously. By **many**, we're referring to a few thousand events, not millions. At this scale, the system's performance bottlenecks become much more apparent.

There are two references that highlight performance issues when inserting a large number of events, each tested under different hardware configurations:

- [Insert 5000 rows per second using PostgreSQL Copy](https://radityaqb.medium.com/insert-5000-rows-per-second-using-postgresql-copy-30fcff1e8fd)
- [Testing Postgres Ingest: INSERT vs. Batch INSERT vs. COPY](https://www.timescale.com/learn/testing-postgres-ingest-insert-vs-batch-insert-vs-copy)

Both references show that while _INSERT_ and _BULK INSERT_ methods perform adequately with small numbers of events, they struggle with larger batches. In contrast, the _COPY_ method offers significantly better performance for high-volume inserts. As the number of events increases, the performance gap between _INSERT_, _BULK INSERT_, and _COPY_ becomes more evident.

:::info

**Why COPY command is better**

According to official PostgreSQL documentation, using the COPY command (you can check it at [Use COPY](https://www.postgresql.org/docs/current/populate.html#POPULATE-COPY-FROM)) is almost always faster than using INSERT, even when PREPARE is used, and multiple insertions are batched within a single transaction.

**Note that loading a large number of rows using COPY is almost always faster than using INSERT, even if PREPARE is used and multiple insertions are batched into a single transaction.**

:::

In the KanthorQ system, we use the _COPY_ command by default for inserting events, due to its superior performance in handling larger batches.

However, there is one caveat: KanthorQ uses the `COPY %s ( %s ) FROM STDIN BINARY` statement under the hood to insert events. This means that events are buffered in memory before being copied into the database, leading to increased memory overhead if the events are large.
