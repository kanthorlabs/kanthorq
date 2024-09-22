---
title: "Limits"
sidebar_label: "Limits"
sidebar_position: 4
---

This page lists all the limits and constraints of the KanthorQ system that you should be aware of.

## Events

- The `body` property of an event can store up to **1GB** of binary data ([Storing Binary Data](https://jdbc.postgresql.org/documentation/binary-data/)). However, we do not recommend storing such large amounts of data due to the potential performance penalties.

## PostgreSQL

- Automatic Prepared Statement Caching feature (mode `QueryExecModeCache`) is incompatible with PgBouncer so that you need to use `default_query_exec_mode` in connection string instead. Example: `postgres://postgres:changemenow@localhost:6432/postgres?sslmode=disable&default_query_exec_mode=exec`. More discussion at [Automatic Prepared Statement Caching](https://github.com/jackc/pgx/wiki/Automatic-Prepared-Statement-Caching) and [Expected good QueryExecMode configurations with PgBouncer 1.21.0?](https://github.com/jackc/pgx/discussions/1784)
