---
title: "Retry and Discard Workflow"
sidebar_label: "Retry and Discard Workflow"
sidebar_position: 1
---

> Anything that can go wrong will go wrong - Murphy's law

When something went wrong at the [Default Workflow](./001-default-workflow.md), your tasks will be marked as `Retryable` first. We need another workflow to handle those `Retryable` tasks if you want to retry it.

The retry process will not be run forever, it will be limited by `attempt_max` option of the consumer which will be set when you initialize the subscriber.

```mermaid
---
title: Retry and Discard Workflow
---
stateDiagram-v2
    direction LR
    [*] --> Available
    Available --> Running
    Running --> Retryable

    Retryable --> Retryable: if attempt_count <= attempt_max
    Retryable --> Discarded: if attempt_count > attempt_max

    Discarded--> [*]
```
