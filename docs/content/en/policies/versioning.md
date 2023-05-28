---
title: "Versioning Policy"
description: "The Authelia Versioning Policy which is important reading for administrators"
date: 2022-12-21T20:48:14+11:00
draft: false
images: []
aliases:
  - /versioning-policy
  - /versioning
---

The __Authelia__ team aims to abide by the [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html) policy. This
means that we use the format `<major>.<minor>.<patch>` for our version numbers, where a change to `major` denotes a
breaking change which will likely require user interaction to upgrade, `minor` which denotes a new feature, and `patch`
denotes a fix.

It is therefore recommended users do not automatically upgrade the `minor` version without reading the patch notes, and
it's critically important users do not upgrade the `major` version without reading the patch notes. You should pin your
version to `4.37` for example to prevent automatic upgrades of the `minor` version, or pin your version to `4` to
prevent automatic upgrade of the `major` version.

We generally do not recommend automated upgrades of critical systems but instead recommend ensuring you are notified an
upgrade exists.

## Supported Versions

The following information is indicative of our support policy:

- We provide support to user questions for 3 `minor` versions at minimum
- We provide bug fixes (as a `patch`) to the latest `minor` version
- We provide vulnerability fixes:
  - As workarounds in the [security advisory](https://github.com/authelia/authelia/security/advisories) (if possible)
  - As patches in the [security advisory](https://github.com/authelia/authelia/security/advisories)
  - To the last 3 `minor` versions upon request

## Major Version Zero

A major version of `v0.x.x` indicates as per the [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html) policy
that there may be breaking changes without warning. Some [components](#components) will be released under this version
while they're in early development.

## Components

Several components may exist at various times. We aim to abide by this policy for all components related to Authelia.
It is important to note that each component has its own version, for example the primary Authelia binary version may be
v4.40.0 but another component such as the [Helm Chart](https://charts.authelia.com) version may be v0.9.0.

This means that a breaking change may occur to one but not the other as these components do not share a version.

## Exceptions

There are exceptions to this versioning policy.

### Advanced Customizations

Some advanced customizations are not guaranteed by the versioning policy. These features require the administrator to
ensure they keep up to date with the changes relevant to their version. While the customizations exist as a feature we
cannot allow these customizations to hinder the development process.

Notable Advanced Customizations:

- Templates:
  - Email
  - Content Security Policy header
- Localization / Internationalization Assets

### Breaking Changes

All features which are marked as:

- beta
- experimental

Notable examples:

- OpenID Connect 1.0
- File Filters

The reasoning is as we develop these features there may be mistakes and we may need to make a change that would normally
be considered a breaking change. As these features graduate from their status to generally available they will move into
our standard versioning policy and lose their exception status.
