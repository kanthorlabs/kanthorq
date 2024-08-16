---
title: "Stream"
sidebar_label: "Stream"
sidebar_position: 4
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

Stream is a persistent, append-only event group that serves specific purposes. For example, you can create a stream with name `order_update` to put all events that are relates to your order into that stream.

```mermaid
---
title: Stream
---
flowchart TB
  Order[Order Service] -- order.created ---> order_update[(kanthorq_stream_order_update)]
  Order[Order Service] -- order.cancelled ---> order_update[(kanthorq_stream_order_update)]
  3PL[3PL Service] -- parcel.lost --> order_update[(kanthorq_stream_order_update)]
```

There are some characteristics of a stream you should know

- An event stays forever in a stream until you explicitly remove it or a stream is deleted (also must be explicit confirmation)
- An event could be read and processed by multiple process (we call it `Consumer`) from the stream and nothing else the event data itself is stored in stream.
- Events in a stream could be only paginated by the order of `event.id` or the tuple of `(event.subject, event.id)`

## Manage streams

When you create or register a stream for you usage, its information will be store in a registry then KanthorQ creates an acutal stream from the for you to store events from the returning registry inforamtion.

```mermaid
---
title: Stream Register Flow
---
sequenceDiagram
  Client ->> +Stream Registry: name: order_update
  Stream Registry -->> -Client: Stream(name: order_update)

  Client ->> +PostgreSQL: kanthorq_stream_order_update
  PostgreSQL -->> -Client: OK
```

### Stream Registry

There is the definition of the `Stream Registry` in different places in KanthorQ

<Tabs>
  <TabItem value="go" label="Go" default>
    ```go
    type StreamRegistry struct {
      Id        string `json:"id" validate:"required"`
      Name      string `json:"name" validate:"required,is_collection_name"`
      CreatedAt int64  `json:"created_at"`
      UpdatedAt int64  `json:"updated_at"`
    }
    ```
  </TabItem>
  <TabItem value="postgresql" label="PostgreSQL">
    ```sql
    TABLE kanthorq_stream_registry (
      id VARCHAR(64) NOT NULL,
      name VARCHAR(256) NOT NULL,
      created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
      updated_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
      PRIMARY KEY (id)
    );
    ```
  </TabItem>
</Tabs>

### Stream

As the definition said about the `Stream`, it's just a **append-only event group** so its definition is a shape of of the `Event`.

```sql
TABLE kanthorq_stream_order_update (
  id VARCHAR(64) NOT NULL,
	subject VARCHAR(256) NOT NULL,
  body BYTEA NOT NULL,
	metadata jsonb NOT NULL,
	created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM CURRENT_TIMESTAMP) * 1000,
	PRIMARY KEY (id)
)
```
