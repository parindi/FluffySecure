---
title: "Skipper"
description: "An integration guide for Authelia and the Skipper reverse proxy"
lead: "A guide on integrating Authelia with the Skipper reverse proxy."
date: 2022-06-15T17:51:47+10:00
draft: false
images: []
menu:
  integration:
    parent: "proxies"
weight: 360
toc: true
aliases:
  - /i/skipper
---

[Skipper] is probably supported by __Authelia__.

*__Important:__ When using these guides it's important to recognize that we cannot provide a guide for every possible
method of deploying a proxy. These guides show a suggested setup only and you need to understand the proxy
configuration and customize it to your needs. To-that-end we include links to the official proxy documentation
throughout this documentation and in the [See Also](#see-also) section.*

## UNDER CONSTRUCTION

It's currently not certain, but fairly likely that [Skipper] is supported by __Authelia__. We wish to add documentation
and thus if anyone has this working please let us know.

We will aim to perform documentation for this on our own but there is no current timeframe.

## Get Started

It's __*strongly recommended*__ that users setting up *Authelia* for the first time take a look at our
[Get Started](../prologue/get-started.md) guide. This takes you through various steps which are essential to
bootstrapping *Authelia*.

## Trusted Proxies

*__Important:__ You should read the [Forwarded Headers] section and this section as part of any proxy configuration.
Especially if you have never read it before.*

*__Important:__ The included example is __NOT__ meant for production use. It's used expressly as an example to showcase
how you can configure multiple IP ranges. You should customize this example to fit your specific architecture and needs.
You should only include the specific IP address ranges of the trusted proxies within your architecture and should not
trust entire subnets unless that subnet only has trusted proxies and no other services.*

## Assumptions and Adaptation

This guide makes a few assumptions. These assumptions may require adaptation in more advanced and complex scenarios. We
can not reasonably have examples for every advanced configuration option that exists. The
following are the assumptions we make:

* Deployment Scenario:
  * Single Host
  * Authelia is deployed as a Container with the container name `authelia` on port `9091`
  * Proxy is deployed as a Container on a network shared with Authelia
* The above assumption means that AUthelia should be accesible to the proxy on `http://authelia:9091` and as such:
  * You will have to adapt all instances of the above URL to be `https://` if Authelia configuration has a TLS key and
    certificate defined
  * You will have to adapt all instances of `authelia` in the URL if:
    * you're using a different container name
    * you deployed the proxy to a different location
  * You will have to adapt all instances of `9091` in the URL if:
    * you have adjusted the default port in the configuration
  * You will have to adapt the entire URL if:
    * Authelia is on a different host to the proxy
* All services are part of the `example.com` domain:
  * This domain and the subdomains will have to be adapted in all examples to match your specific domains unless you're
    just testing or you want ot use that specific domain

## Potential

Support for [Skipper] should be possible via [Skipper]'s
[Webhook Filter](https://opensource.zalando.com/skipper/reference/filters/#webhook).

## See Also

* [Skipper Webhook Filter Documentation](https://opensource.zalando.com/skipper/reference/filters/#webhook)
* [Forwarded Headers]

[Skipper]: https://opensource.zalando.com/skipper/
[Forwarded Headers]: fowarded-headers
