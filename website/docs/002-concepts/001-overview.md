---
title: "Overview"
sidebar_label: "Overview"
sidebar_position: 1
---

Let's discover what the KanthorQ architecture is and how the components communicate with each other.

## Architecture

The KanthorQ architecture, depicted in the diagram below, comprises four components:

- **Publisher**: Application code responsible for inserting events into the KanthorQ system. This can be a Command Line Tool (which I'm maintaining) or Golang code within your application.
- **Stream**: Receives events and persists them within the KanthorQ system, categorized by topics.
- **Consumer**: Stores tasks generated from events in the Stream. One event can generate many tasks in different _Consumers_, but in the same _Consumer_, there must be only one task that belongs to an event
- **Subscriber**: Part of your application that pulls tasks to execute your business logic.

![KanthorQ Architecture](/002-concepts/001-overview/kanthorq-architecture.svg)
