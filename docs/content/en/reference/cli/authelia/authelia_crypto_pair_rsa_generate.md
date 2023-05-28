---
title: "authelia crypto pair rsa generate"
description: "Reference for the authelia crypto pair rsa generate command."
lead: ""
date: 2022-06-27T18:27:57+10:00
draft: false
images: []
menu:
  reference:
    parent: "cli-authelia"
weight: 905
toc: true
---

## authelia crypto pair rsa generate

Generate a cryptographic RSA key pair

### Synopsis

Generate a cryptographic RSA key pair.

This subcommand allows generating an RSA key pair.

```
authelia crypto pair rsa generate [flags]
```

### Examples

```
authelia crypto pair rsa generate --help
```

### Options

```
  -b, --bits int                  number of RSA bits for the certificate (default 2048)
  -d, --directory string          directory where the generated keys, certificates, etc will be stored
      --file.private-key string   name of the file to export the private key data to (default "private.pem")
      --file.public-key string    name of the file to export the public key data to (default "public.pem")
  -h, --help                      help for generate
      --pkcs8                     force PKCS #8 ASN.1 format
```

### Options inherited from parent commands

```
  -c, --config strings                        configuration files or directories to load, for more information run 'authelia -h authelia config' (default [configuration.yml])
      --config.experimental.filters strings   list of filters to apply to all configuration files, for more information run 'authelia -h authelia filters'
```

### SEE ALSO

* [authelia crypto pair rsa](authelia_crypto_pair_rsa.md)	 - Perform RSA key pair cryptographic operations

