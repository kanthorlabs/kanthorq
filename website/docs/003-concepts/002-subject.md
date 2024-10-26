---
title: "Subject"
sidebar_label: "Subject"
sidebar_position: 2
---

A Subject in the KanthorQ system is a hierarchical entity that categorizes events. Understanding how subjects are structured and used in KanthorQ is essential for leveraging them effectively in your workflows.

## Recommendation subject design decisions

- A Subject in KanthorQ follows a left-to-right naming hierarchy, moving from general to specific. Here are examples illustrating this hierarchy:

  - `ecommerce.order.created` – Simple name definition.
  - `tier_standard.order.created` – Tier-specific subject.
  - `internal.testing.order.created` – Naming for internal testing environments.
  - **Environment-specific**: Separate operational environments, such as ecommerce.`order.created.prod` and `ecommerce.order.created.uat`, to facilitate testing or debugging.
  - **Versioned subjects**: Use versions like `ecommerce.order.created.v1` and `ecommerce.order.created.v2` for tracking updates over time.

- **Abstract Layer**: Subjects should represent business logic, not technical details.
- **Longevity**: Subjects should be designed to last throughout the business's lifecycle. If changes are required, consider creating a new version or an entirely new subject to replace outdated ones.

## Define your subject

KanthorQ's subject definition aligns with [NATS Subject-Based Messaging](https://docs.nats.io/nats-concepts/subjects), supporting an interest-based messaging system. Here are the rules for defining subjects:

- A `subject` is a case-sensitive string composed of multiple tokens separated by . (dot).
- Each `token` is a string, allowing any character except `.`, `*`, and `>`.

There are some recommendations from the section [Characters allowed and recommended for subject names](https://docs.nats.io/nats-concepts/subjects#characters-allowed-and-recommended-for-subject-names) I would recommend when you want to design your subject

- **Allowed characters**: Any Unicode character except `null`, `space`, `.`, `*` and `>`
- **Recommended characters**: (`a` - `z`), (`A` - `Z`), (`0` - `9`), `-` and `_` (names are case sensitive, and cannot contain whitespace).
- **Naming Conventions**: If you want to delimit words, use either CamelCase as in `MyServiceOrderCreate` or - and `_` as in `my-service-order-create`
- **Special characters**: The period `.` (which is used to separate the tokens in the subject) and `*` and also `>` (the `*` and `>` are used as wildcards) are reserved and cannot be used.

## Matching your subject

In KanthorQ, you can match subjects to patterns using three methods:

- **Exact Match**: Matches subjects with exact case-sensitive text.
- Single-Token Wildcard (`*`): Matches exactly one token in the subject. For example, `order.*` matches `order.created` or `order.updated` but does not match `order`, `order.created.v1`, or `order.updated.v1`.
- Multi-Token Wildcard (`>`): Using `>` wildcard to match multiple tokens in a subject. For example, `order.>` matches `order.created`, `order.updated`, `order.created.v1`, and `order.updated.v1` but does not match `order`.

For more pattern matching examples, refer to the [TestMatchSubject](https://github.com/kanthorlabs/kanthorq/blob/main/core/utils_test.go) test.

## Final Thoughts

Nothing, go test and find out some gochas ;D
