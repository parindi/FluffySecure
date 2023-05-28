---
title: "authelia-scripts suites teardown"
description: "Reference for the authelia-scripts suites teardown command."
lead: ""
date: 2022-06-15T17:51:47+10:00
draft: false
images: []
menu:
  reference:
    parent: "cli-authelia-scripts"
weight: 925
toc: true
---

## authelia-scripts suites teardown

Teardown a test suite environment

### Synopsis

Teardown a test suite environment.

Suites can be listed with the authelia-scripts suites list command.

```
authelia-scripts suites teardown [suite] [flags]
```

### Examples

```
authelia-scripts suites setup Standalone
```

### Options

```
  -h, --help   help for teardown
```

### Options inherited from parent commands

```
      --buildkite          Set CI flag for Buildkite
      --log-level string   Set the log level for the command (default "info")
```

### SEE ALSO

* [authelia-scripts suites](authelia-scripts_suites.md)	 - Commands related to suites management

