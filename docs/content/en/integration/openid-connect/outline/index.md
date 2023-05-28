---
title: "Outline"
description: "Integrating Outline with the Authelia OpenID Connect 1.0 Provider."
lead: ""
date: 2022-08-12T09:11:42+10:00
draft: false
images: []
menu:
  integration:
    parent: "openid-connect"
weight: 620
toc: true
community: true
---

## Tested Versions

* [Authelia]
  * [v4.36.4](https://github.com/authelia/authelia/releases/tag/v4.36.4)
* [Outline]
  * 0.65.2

## Before You Begin

{{% oidc-common %}}

### Assumptions

This example makes the following assumptions:

* __Application Root URL:__ `https://outline.example.com`
* __Authelia Root URL:__ `https://auth.example.com`
* __Client ID:__ `outline`
* __Client Secret:__ `insecure_secret`

*__Important Note:__ At the time of this writing [Outline] requires the `offline_access` scope by default. Failure to include this scope will result
in an error as [Outline] will attempt to use a refresh token that is never issued.*

## Configuration

### Application

To configure [Outline] to utilize Authelia as an [OpenID Connect 1.0] Provider:

1. Configure the following environment options:
```text
URL=https://outline.example.com
FORCE_HTTPS=true

OIDC_CLIENT_ID=outline
OIDC_CLIENT_SECRET=insecure_secret
OIDC_AUTH_URI=https://auth.example.com/api/oidc/authorization
OIDC_TOKEN_URI=https://auth.example.com/api/oidc/token
OIDC_USERINFO_URI=https://auth.example.com/api/oidc/userinfo
OIDC_USERNAME_CLAIM=preferred_username
OIDC_DISPLAY_NAME=Authelia
OIDC_SCOPES="openid offline_access profile email"
```

### Authelia

The following YAML configuration is an example __Authelia__
[client configuration](../../../configuration/identity-providers/openid-connect/clients.md) for use with [Outline]
which will operate with the above example:

```yaml
identity_providers:
  oidc:
    ## The other portions of the mandatory OpenID Connect 1.0 configuration go here.
    ## See: https://www.authelia.com/c/oidc
    clients:
    - id: 'outline'
      description: 'Outline'
      secret: '$pbkdf2-sha512$310000$c8p78n7pUMln0jzvd4aK4Q$JNRBzwAo0ek5qKn50cFzzvE9RXV88h1wJn5KGiHrD0YKtZaR/nCb2CJPOsKaPK0hjf.9yHxzQGZziziccp6Yng'  # The digest of 'insecure_secret'.
      public: false
      authorization_policy: 'two_factor'
      redirect_uris:
        - 'https://outline.example.com/auth/oidc.callback'
      scopes:
        - 'openid'
        - 'offline_access'
        - 'profile'
        - 'email'
      userinfo_signing_alg: 'none'
```

## See Also

* [Outline OpenID Connect Documentation](https://app.getoutline.com/share/770a97da-13e5-401e-9f8a-37949c19f97e/doc/oidc-8CPBm6uC0I)

[Authelia]: https://www.authelia.com
[Outline]: https://www.getoutline.com/
[OpenID Connect 1.0]: ../../openid-connect/introduction.md
