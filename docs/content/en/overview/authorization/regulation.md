---
title: "Regulation"
description: "Regulation of failed attempts is an important function of an IAM system."
lead: "Regulation of failed attempts is an important function of an IAM system."
date: 2022-06-15T17:51:47+10:00
draft: false
images: []
menu:
  overview:
    parent: "authorization"
weight: 320
toc: false
aliases:
  - /docs/features/regulation.html
---

__Authelia__ takes the security of users very seriously and comes with a way to avoid brute-forcing the first factor
credentials by regulating the authentication attempts and temporarily banning an account when too many attempts have
been made.

## Configuration

Please check the dedicated [documentation](../../configuration/security/regulation.md).
