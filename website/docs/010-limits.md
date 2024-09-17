---
title: "Limits"
sidebar_label: "Limits"
sidebar_position: 4
---

This page lists all the limits and constraints of the KanthorQ system that you should be aware of.

## Events

- The `body` property of an event can store up to **1GB** of binary data ([Storing Binary Data](https://jdbc.postgresql.org/documentation/binary-data/)). However, we do not recommend storing such large amounts of data due to the potential performance penalties.
