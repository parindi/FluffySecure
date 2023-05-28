---
title: "Secrets"
description: "Using the Secrets Configuration Method."
lead: "Authelia allows providing configuration via secrets method. This section describes how to implement this."
date: 2020-02-29T01:43:59+01:00
draft: false
images: []
menu:
  configuration:
    parent: "methods"
weight: 101400
toc: true
aliases:
  - /c/secrets
  - /docs/configuration/secrets.html
---

Configuration of *Authelia* requires several secrets and passwords. Even if they can be set in the configuration file or
standard environment variables, the recommended way to set secrets is to use this configuration method as described below.

See the [security](#security) section for more information.

## Layers

*__Important Note:__* While this method is the third layer of the layered configuration model as described by the
[introduction](introduction.md#layers), this layer is special in as much as *Authelia* will not start if you define
a secret as well as any other configuration method.

For example if you define `jwt_secret` in the [files method](files.md) and/or `AUTHELIA_JWT_SECRET` in the
[environment method](environment.md), as well as the `AUTHELIA_JWT_SECRET_FILE`, this will cause the aforementioned error.

## Security

This method is a slight improvement over the security of the other methods as it allows you to easily separate your
configuration in a logically secure way.

## Environment variables

A secret value can be loaded by *Authelia* when the configuration key ends with one of the following words: `key`,
`secret`, `password`, or `token`.

If you take the expected environment variable for the configuration option with the `_FILE` suffix at the end. The value
of these environment variables must be the path of a file that is readable by the Authelia process, if they are not,
*Authelia* will fail to load. Authelia will automatically remove the newlines from the end of the files contents.

For instance the LDAP password can be defined in the configuration
at the path __authentication_backend.ldap.password__, so this password
could alternatively be set using the environment variable called
__AUTHELIA_AUTHENTICATION_BACKEND_LDAP_PASSWORD_FILE__.

Here is the list of the environment variables which are considered secrets and can be defined. Please note that only
secrets can be loaded into the configuration if they end with one of the suffixes above, you can set the value of any
other configuration using the environment but instead of loading a file the value of the environment variable is used.

{{% table-config-keys secrets="true" %}}

[server.tls.key]: ../miscellaneous/server.md#key
[jwt_secret]: ../miscellaneous/introduction.md#jwtsecret
[duo_api.integration_key]: ../second-factor/duo.md#integrationkey
[duo_api.secret_key]: ../second-factor/duo.md#secretkey
[session.secret]: ../session/introduction.md#secret
[session.redis.password]: ../session/redis.md#password
[session.redis.tls.certificate_chain]: ../session/redis.md#tls
[session.redis.tls.private_key]: ../session/redis.md#tls
[session.redis.high_availability.sentinel_password]: ../session/redis.md#sentinelpassword
[storage.encryption_key]: ../storage/introduction.md#encryptionkey
[storage.mysql.password]: ../storage/mysql.md#password
[storage.mysql.tls.certificate_chain]: ../storage/mysql.md#tls
[storage.mysql.tls.private_key]: ../storage/mysql.md#tls
[storage.postgres.password]: ../storage/postgres.md#password
[storage.postgres.tls.certificate_chain]: ../storage/postgres.md#tls
[storage.postgres.tls.private_key]: ../storage/postgres.md#tls
[storage.postgres.ssl.key]: ../storage/postgres.md
[notifier.smtp.password]: ../notifications/smtp.md#password
[notifier.smtp.tls.certificate_chain]: ../notifications/smtp.md#tls
[notifier.smtp.tls.private_key]: ../notifications/smtp.md#tls
[authentication_backend.ldap.password]: ../first-factor/ldap.md#password
[authentication_backend.ldap.tls.certificate_chain]: ../first-factor/ldap.md#tls
[authentication_backend.ldap.tls.private_key]: ../first-factor/ldap.md#tls
[identity_providers.oidc.issuer_certificate_chain]: ../identity-providers/openid-connect.md#issuercertificatechain
[identity_providers.oidc.issuer_private_key]: ../identity-providers/openid-connect.md#issuerprivatekey
[identity_providers.oidc.hmac_secret]: ../identity-providers/openid-connect.md#hmacsecret


## Secrets in configuration file

If for some reason you decide on keeping the secrets in the configuration file, it is strongly recommended that you
ensure the permissions of the configuration file are appropriately set so that other users or processes cannot access
this file. Generally the UNIX permissions that are appropriate are 0600.

## Secrets exposed in an environment variable

In all versions 4.30.0+ you can technically set secrets using the environment variables without the `_FILE` suffix by
setting the value to the value you wish to set in configuration, however we strongly urge people not to use this option
and instead use the file-based secrets above.

Prior to implementing file secrets the only way you were able to define secret values was either via configuration or
via environment variables in plain text.

See [this article](https://diogomonica.com/2017/03/27/why-you-shouldnt-use-env-variables-for-secret-data/) for reasons
why setting them via the file counterparts is highly encouraged.

## Examples

See the [Docker Integration](../../integration/deployment/docker.md) and
[Kubernetes Integration](../../integration/kubernetes/secrets.md) guides for examples of secrets.
