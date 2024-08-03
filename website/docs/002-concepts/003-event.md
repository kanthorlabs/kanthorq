---
title: "Event"
sidebar_label: "Event"
sidebar_position: 3
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

`Event` is an entity that represents a data transfer object (DTO) between publishers and streams in KanthorQ. It is similar with HTTP request in client-server commication over HTTP or proto definition in gRPC.

There is the definition of the `Event` in different places in KanthorQ

<Tabs>
  <TabItem value="go" label="Go" default>
    ```go
    type Event struct {
      Id        string         `json:"id"`
      Subject     string         `json:"subject"`
      Body      []byte         `json:"body"`
      Metadata  map[string]any `json:"metadata"`
      CreatedAt int64          `json:"created_at"`
    }
    ```
  </TabItem>
  <TabItem value="postgresql" label="PostgreSQL">
    ```sql
    TABLE kanthorq_stream_order_update (
      id VARCHAR(64) NOT NULL,
      subject VARCHAR(128) NOT NULL,
      body BYTEA NOT NULL,
      metadata jsonb NOT NULL,
      created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000
    )
    ```
  </TabItem>
</Tabs>

:::info

You can be confused why we use the name `kanthorq_stream_order_update` to represent the `Event` struct in PostgreSQL. Check [Stream Concept](/docs/concepts/stream#stream) for more information.

:::

To make the communication is easy to integrate we have define some characteristics you MUST and SHOULD follow.

MUST follow requirements includes:

- The `id` property must be Lexicographically Sortable Identifier. For example if you want to use UUID, please use UUIDv7. We are recommend you use the [ULID](https://github.com/ulid/spec)
- The `subject` property must multiple alphanumeric+hypen strings that could be seperate by a dot.

  - OK: `order`, `order.created`, `v2.order.created`, `058434268238.order.created`, `-1002223543143.subscription.created`, `ap-southeast-1.user.created`
  - KO: `order.`, `order.*`

- The `body` property stores arbitrary bytes so only the sender and the receiver know what it stores. It's usefull in the end-to-end encryption implementation where the sender encrypt the data before sending. Then only the receiver know how to descrypt it.

SHOULD follow requirements includes:

- The `metadata` property is an json object that store additional information about the event. You can use it to implement event filter logic or telemetry carier for example.
- The `created_at` property should contains the timestamp that was includes in the `id` property. That allow you get events in datatime order just by ordering events by `id`
