---
title: "Event"
sidebar_label: "Event"
sidebar_position: 3
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

An `Event` in KanthorQ represents a data transfer object (DTO) used to communicate between publishers and streams.

<Tabs>
  <TabItem value="go" label="Go" default>
    ```go
    type Event struct {
      Id        string   `json:"id" validate:"required"`
      Subject   string   `json:"subject" validate:"required,is_subject"`
      Body      []byte   `json:"body" validate:"required"`
      Metadata  Metadata `json:"metadata" validate:"required"`
      CreatedAt int64    `json:"created_at"`
    }
    ```
  </TabItem>
  <TabItem value="postgresql" label="PostgreSQL">
    ```sql
    -- Example of Event structure in a Stream
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

You might wonder why the Event struct is named `kanthorq_stream_order_update` in PostgreSQL—this is related to the Stream concept. Check the [Stream Concept](./004-stream.md#stream) section for more information.

:::

To ensure seamless communication and integration, we've established some rules that you MUST and SHOULD follow:

**MUST Follow Requirements:**

- The `id` property must be a Lexicographically Sortable Identifier. For example, if you're using UUIDs, please use UUIDv7. However, we recommend using [ULID](https://github.com/ulid/spec) for this purpose.
- The `subject` property must consist of multiple alphanumeric and hyphenated strings, separated by dots.

  - Correct Examples: `order`, `order.created`, `v2.order.created`, `058434268238.order.created`, `-1002223543143.subscription.created`, `ap-southeast-1.user.created`
  - Incorrect Examples: `order.`, `order.-`, `order..`

- The `body` property stores arbitrary bytes, which only the sender and receiver can interpret. This is particularly useful for implementing end-to-end encryption, where the sender encrypts the data before sending it, and only the receiver can decrypt it.

**SHOULD Follow Requirements:**

- The `metadata` property should be a JSON object that stores additional information about the event. This can be used for implementing event filtering logic or as a telemetry carrier.
- The `created_at` property should contain the timestamp included in the id property. This ensures that events can be sorted chronologically by simply ordering them by their id.
