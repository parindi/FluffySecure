---
title: "Metrics"
description: "Configuring the Metrics Telemetry settings"
lead: "Configuring the Metrics Telemetry settings."
date: 2022-06-15T17:51:47+10:00
draft: false
images: []
menu:
  configuration:
    parent: "telemetry"
weight: 108200
toc: true
---

*Authelia* allows administrators to configure a [Prometheus] Metrics Exporter.

## Configuration

{{< config-alert-example >}}

```yaml
telemetry:
  metrics:
    enabled: false
    address: 'tcp://:9959'
    buffers:
      read: 4096
      write: 4096
    timeouts:
      read: '6s'
      write: '6s'
      idle: '30s'
```

## Options

This section describes the individual configuration options.

### enabled

{{< confkey type="boolean" default="false" required="no" >}}

Determines if the [Prometheus] HTTP Metrics Exporter is enabled.

### address

{{< confkey type="address" default="tcp://:9959" required="no" >}}

*__Reference Note:__ This configuration option uses the [address common syntax](../prologue/common.md#address). Please
see the [documentation](../prologue/common.md#address) on this format for more information.*

Configures the listener address for the [Prometheus] Metrics Exporter HTTP Server. The address itself is a listener and
the scheme must either be the `unix` scheme or one of the `tcp` schemes.

### buffers

*__Reference Note:__ This configuration option uses the
[Server buffers common structure](../prologue/common.md#server-buffers). Please see the
[documentation](../prologue/common.md#server-buffers) on this structure for more information.*

Configures the server buffers.

### timeouts

*__Reference Note:__ This configuration option uses the
[Server timeouts common structure](../prologue/common.md#server-timeouts). Please see the
[documentation](../prologue/common.md#server-timeouts) on this structure for more information.*

Configures the server timeouts.

## See More

- [Telemetry Reference Documentation](../../reference/guides/metrics.md)

[Prometheus]: https://prometheus.io/
