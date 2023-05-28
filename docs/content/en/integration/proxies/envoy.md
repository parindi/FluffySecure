---
title: "Envoy"
description: "An integration guide for Authelia and the Envoy reverse proxy"
lead: "A guide on integrating Authelia with the Envoy reverse proxy."
date: 2022-06-15T17:51:47+10:00
draft: false
images: []
menu:
  integration:
    parent: "proxies"
weight: 330
toc: true
aliases:
  - /i/envoy
---

[Envoy] is supported by __Authelia__.

*__Important:__ When using these guides it's important to recognize that we cannot provide a guide for every possible
method of deploying a proxy. These guides show a suggested setup only and you need to understand the proxy
configuration and customize it to your needs. To-that-end we include links to the official proxy documentation
throughout this documentation and in the [See Also](#see-also) section.*

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

## Implementation

[Envoy] utilizes the [ExtAuthz](../../reference/guides/proxy-authorization.md#extauthz) Authz implementation. The
associated [Metadata](../../reference/guides/proxy-authorization.md#extauthz-metadata) should be considered required.

The examples below assume you are using the default
[Authz Endpoints Configuration](../../configuration/miscellaneous/server-endpoints-authz.md) or one similar to the
following minimal configuration:

```yaml
server:
  endpoints:
    authz:
      ext-authz:
        implementation: 'ExtAuthz'
```

## Configuration

Below you will find commented examples of the following configuration:

* Authelia Portal
* Protected Endpoint (Nextcloud)

### Example

Support for [Envoy] is possible with Authelia v4.37.0 and higher via the [Envoy] proxy [external authorization] filter.

[external authorization]: https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/ext_authz/v3/ext_authz.proto.html#extensions-filters-http-ext-authz-v3-extauthz

{{< details "docker-compose.yaml" >}}
```yaml
---
version: "3.8"
networks:
  net:
    driver: 'bridge'
services:
  envoy:
    container_name: 'envoy'
    image: 'envoyproxy/envoy:v1.24'
    restart: 'unless-stopped'
    networks:
      net: {}
    ports:
      - '80:8080'
      - '443:8443'
    volumes:
      - '${PWD}/data/envoy/envoy.yaml:/etc/envoy/envoy.yaml:ro'
      - '${PWD}/data/certificates:/certificates:ro'
  authelia:
    container_name: 'authelia'
    image: 'authelia/authelia'
    restart: 'unless-stopped'
    networks:
      net: {}
    expose:
      - 9091
    volumes:
      - '${PWD}/data/authelia/config:/config'
    environment:
      TZ: 'Australia/Melbourne'
  nextcloud:
    container_name: 'nextcloud'
    image: 'linuxserver/nextcloud'
    restart: 'unless-stopped'
    networks:
      net: {}
    expose:
      - 443
    volumes:
      - '${PWD}/data/nextcloud/config:/config'
      - '${PWD}/data/nextcloud/data:/data'
    environment:
      PUID: '1000'
      PGID: '1000'
      TZ: 'Australia/Melbourne'
```
{{< /details >}}

{{< details "envoy.yaml" >}}
```yaml
static_resources:
  listeners:
    - name: 'listener_http'
      address:
        socket_address:
          address: '0.0.0.0'
          port_value: 8080
      filter_chains:
        - filters:
            - name: 'envoy.filters.network.http_connection_manager'
              typed_config:
                "@type": 'type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager'
                codec_type: 'auto'
                stat_prefix: 'ingress_http'
                route_config:
                  name: 'local_route'
                  virtual_hosts:
                    - name: 'backend'
                      domains: ['*']
                      routes:
                        - match:
                            prefix: '/'
                          redirect:
                            https_redirect: true
                  http_filters:
                    - name: 'envoy.filters.http.router'
                      typed_config:
                        "@type": 'type.googleapis.com/envoy.extensions.filters.http.router.v3.Router'
    - name: 'listener_https'
      address:
        socket_address:
          address: '0.0.0.0'
          port_value: 8443
      filter_chains:
        - filters:
            - name: 'envoy.filters.network.http_connection_manager'
              typed_config:
                "@type": 'type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager'
                stat_prefix: 'ingress_http'
                use_remote_address: true
                skip_xff_append: false
                route_config:
                  name: 'local_route'
                  virtual_hosts:
                    - name: 'whoami_service'
                      domains: ["nextcloud.example.com"]
                      routes:
                        - match:
                            prefix: "/"
                          route:
                            cluster: 'nextcloud'
                    - name: 'authelia_service'
                      domains: ['auth.example.com']
                      typed_per_filter_config:
                        envoy.filters.http.ext_authz:
                          "@type": 'type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthzPerRoute'
                          disabled: true
                      routes:
                        - match:
                            prefix: '/'
                          route:
                            cluster: 'authelia'
                http_filters:
                  - name: 'envoy.filters.http.ext_authz'
                    typed_config:
                      "@type": 'type.googleapis.com/envoy.extensions.filters.http.ext_authz.v3.ExtAuthz'
                      transport_api_version: 'v3'
                      allowed_headers:
                        patterns:
                          - exact: 'authorization'
                          - exact: 'proxy-authorization'
                          - exact: 'accept'
                          - exact: 'cookie'
                      http_service:
                        path_prefix: '/api/authz/ext-authz/'
                        server_uri:
                          uri: 'authelia:9091'
                          cluster: 'authelia'
                          timeout: '0.25s'
                        authorization_request:
                          allowed_headers:
                            patterns:
                              - exact: 'authorization'
                              - exact: 'proxy-authorization'
                              - exact: 'accept'
                              - exact: 'cookie'
                          headers_to_add:
                            - key: 'X-Forwarded-Proto'
                              value: '%REQ(:SCHEME)%'
                            ## The following commented lines are for configuring the Authelia URL in the proxy. We
                            ## strongly suggest this is configured in the Session Cookies section of the Authelia configuration.
                            # - key: X-Authelia-URL
                            #   value: https://auth.example.com
                        authorization_response:
                          allowed_upstream_headers:
                            patterns:
                              - exact: 'authorization'
                              - exact: 'proxy-authorization'
                              - prefix: 'remote-'
                              - prefix: 'authelia-'
                          allowed_client_headers:
                            patterns:
                              - exact: 'set-cookie'
                          allowed_client_headers_on_success:
                            patterns:
                              - exact: 'set-cookie'
                      failure_mode_allow: false
                  - name: 'envoy.filters.http.router'
                    typed_config:
                      "@type": 'type.googleapis.com/envoy.extensions.filters.http.router.v3.Router'
  clusters:
    - name: 'nextcloud'
      connect_timeout: '0.25s'
      type: 'logical_dns'
      dns_lookup_family: 'v4_only'
      lb_policy: 'round_robin'
      load_assignment:
        cluster_name: 'nextcloud'
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: 'nextcloud'
                      port_value: 80
    - name: 'authelia'
      connect_timeout: '0.25s'
      type: 'logical_dns'
      dns_lookup_family: 'v4_only'
      lb_policy: 'round_robin'
      load_assignment:
        cluster_name: 'authelia'
        endpoints:
          - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                      address: 'authelia'
                      port_value: 9091
layered_runtime:
  layers:
    - name: 'static_layer_0'
      static_layer:
        envoy:
          resource_limits:
            listener:
              example_listener_name:
                connection_limit: 10000
        overload:
          global_downstream_max_connections: 50000
```
{{< /details >}}

## See Also

* [Envoy External Authorization Documentation](https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/ext_authz/v3/ext_authz.proto.html#extensions-filters-http-ext-authz-v3-extauthz)
* [Forwarded Headers]

[Envoy]: https://www.envoyproxy.io/
[Forwarded Headers]: fowarded-headers
