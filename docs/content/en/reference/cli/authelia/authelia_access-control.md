---
title: "authelia access-control"
description: "Reference for the authelia access-control command."
lead: ""
date: 2022-06-15T17:51:47+10:00
draft: false
images: []
menu:
  reference:
    parent: "cli-authelia"
weight: 905
toc: true
---

## authelia access-control

Helpers for the access control system

### Synopsis

Helpers for the access control system.

### Examples

```
authelia access-control --help
```

### Options

```
  -h, --help   help for access-control
```

### Options inherited from parent commands

```
  -c, --config strings                        configuration files or directories to load, for more information run 'authelia -h authelia config' (default [configuration.yml])
      --config.experimental.filters strings   list of filters to apply to all configuration files, for more information run 'authelia -h authelia filters'
```

### SEE ALSO

* [authelia](authelia.md)	 - authelia untagged-unknown-dirty (master, unknown)
* [authelia access-control check-policy](authelia_access-control_check-policy.md)	 - Checks a request against the access control rules to determine what policy would be applied

