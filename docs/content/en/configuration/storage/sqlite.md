---
title: "SQLite3"
description: "SQLite3 Configuration"
lead: "The SQLite3 storage provider."
date: 2022-06-15T17:51:47+10:00
draft: false
images: []
menu:
  configuration:
    parent: "storage"
weight: 106500
toc: true
aliases:
  - /docs/configuration/storage/sqlite.html
---

If you don't have a SQL server, you can use [SQLite](https://en.wikipedia.org/wiki/SQLite).
However please note that this setup will prevent you from running multiple
instances of Authelia since the database will be a local file.

Use of this storage provider leaves Authelia [stateful](../../overview/authorization/statelessness.md). It's important
in highly available scenarios to use one of the other providers, and we highly recommend it in production environments,
but this requires you setup an external database such as [PostgreSQL](postgres.md).

## Configuration

{{< config-alert-example >}}

```yaml
storage:
  encryption_key: 'a_very_important_secret'
  local:
    path: '/config/db.sqlite3'
```

## Options

This section describes the individual configuration options.

### encryption_key

See the [encryption_key docs](introduction.md#encryptionkey).

### path

{{< confkey type="string" required="yes" >}}

The path where the SQLite3 database file will be stored. It will be created if the file does not exist.
