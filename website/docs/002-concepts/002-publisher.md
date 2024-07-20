---
title: "Publisher"
sidebar_label: "Publisher"
sidebar_position: 2
---

The Publisher has only one responsibility, it receives your events, push it into a stream and report back to you whether those events are inserted successfully or not. It does not thing about organizing your data into a time-series shape so make sure the ID of your event is Lexicographically Sortable Identifier
