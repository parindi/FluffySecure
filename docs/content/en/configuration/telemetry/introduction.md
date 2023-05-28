---
title: "Telemetry"
description: "Configuring the Telemetry settings"
lead: "Configuring the Telemetry settings."
date: 2022-06-15T17:51:47+10:00
draft: false
images: []
menu:
  configuration:
    parent: "telemetry"
weight: 108100
toc: true
---

*Authelia* allows collecting telemetry for the purpose of monitoring it. At the present time we only allow collecting
[metrics](metrics.md). These [metrics](metrics.md) are stored in memory and must be scraped manually by the
administrator.

No metrics or telemetry are reported from an *Authelia* binary to any location the administrator doesn't explicitly
configure. This means by default all metrics are disabled.
