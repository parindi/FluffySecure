---
title: "authelia-scripts suites test"
description: "Reference for the authelia-scripts suites test command."
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

## authelia-scripts suites test

Run a test suite

### Synopsis

Run a test suite.

Suites can be listed with the authelia-scripts suites list command.

```
authelia-scripts suites test [suite] [flags]
```

### Examples

```
authelia-scripts suites test Standalone
```

### Options

```
      --failfast      Stops tests on first failure
      --headless      Run tests in headless mode
  -h, --help          help for test
      --test string   The single test to run
```

### Options inherited from parent commands

```
      --buildkite          Set CI flag for Buildkite
      --log-level string   Set the log level for the command (default "info")
```

### SEE ALSO

* [authelia-scripts suites](authelia-scripts_suites.md)	 - Commands related to suites management

