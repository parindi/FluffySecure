---
title: "LDAP"
description: "A reference guide on the LDAP implementation specifics"
lead: "This section contains reference documentation for Authelia's LDAP implementation specifics."
date: 2022-06-17T21:03:47+10:00
draft: false
images: []
menu:
  reference:
    parent: "guides"
weight: 220
toc: true
aliases:
  - /r/ldap
---

## Binding

When it comes to LDAP there are several considerations for deciding how to bind to the LDAP server.

### Unauthenticated Binding

The most insecure method is unauthenticated binds. They are generally considered insecure due to the fact allowing them
at all ensures anyone with any level of network access can easily obtain objects and their attributes.

Authelia does support unauthenticated binds but it is not by default, you must configure the
[permit_unauthenticated_bind](../../configuration/first-factor/ldap.md#permitunauthenticatedbind) configuration
option.

### End-User Binding

One method to bind to the server that is favored by a lot of people is binding to the LDAP server as the end user. While
this is more secure than methods such as [Unauthenticated Binding](#unauthenticated-binding) the drawback is that it can
only be used securely at the time the user enters their credentials. Storing a password in memory in general is not very
secure and prone to breakage due to outside influences (i.e. the user changes their password).

In addition, this method is not compatible with the password reset / forgot password flow at all (not to be confused
with a change password flow).

Authelia doesn't currently support such a binding method excluding for checking user passwords.

### Service-User Binding

This is the most common method of binding to LDAP. This involves setting up a special service user with a complex
password which has the minimum permissions required to do the tasks required.

Authelia primarily supports this method.

## Implementation Guide

The following implementations exist:

- `custom`:
  - Not specific to any particular LDAP provider
- `activedirectory`:
  - Specific configuration defaults for [Active Directory]
  - Special implementation details:
    - Includes a special encoding format required for changing passwords with [Active Directory]
- `rfc2307bis`:
  - Specific configuration defaults for [RFC2307bis]
  - No special implementation details
- `freeipa`:
  - Specific configuration defaults for [FreeIPA]
  - No special implementation details
- `lldap`:
  - Specific configuration defaults for [lldap]
  - No special implementation details
- `glauth`:
  - Specific configuration defaults for [GLAuth]
  - No special implementation details

[Active Directory]: https://learn.microsoft.com/en-us/windows-server/identity/ad-ds/active-directory-domain-services
[FreeIPA]: https://www.freeipa.org/
[lldap]: https://github.com/nitnelave/lldap
[GLAuth]: https://glauth.github.io/
[RFC2307bis]: https://datatracker.ietf.org/doc/html/draft-howard-rfc2307bis-02

### Filter replacements

Various replacements occur in the user and groups filter. The replacements either occur at startup or upon an LDAP
search which is indicated by the phase column.

The phases exist to optimize performance. The replacements in the startup phase are replaced once before the connection
is ever established. In addition to this, during the startup phase we purposefully check the filters for which search
phase replacements exist so we only have to check if the replacement is necessary once, and we don't needlessly perform
every possible replacement on every search regardless of if it's needed or not.

#### Users filter replacements

|       Placeholder        |  Phase  |                                                   Replacement                                                    |
|:------------------------:|:-------:|:----------------------------------------------------------------------------------------------------------------:|
|   {username_attribute}   | startup |                                        The configured username attribute                                         |
|     {mail_attribute}     | startup |                                          The configured mail attribute                                           |
| {display_name_attribute} | startup |                                      The configured display name attribute                                       |
|         {input}          | search  |                                        The input into the username field                                         |
| {date-time:generalized}  | search  |          The current UTC time formatted as a LDAP generalized time in the format of `20060102150405.0Z`          |
|     {date-time:unix}     | search  |                                    The current time formatted as a Unix epoch                                    |
| {date-time:microsoft-nt} | search  | The current time formatted as a Microsoft NT epoch which is used by some Microsoft [Active Directory] attributes |

#### Groups filter replacements

| Placeholder | Phase  |                                Replacement                                |
|:-----------:|:------:|:-------------------------------------------------------------------------:|
|   {input}   | search |                     The input into the username field                     |
| {username}  | search | The username from the profile lookup obtained from the username attribute |
|    {dn}     | search |              The distinguished name from the profile lookup               |

### Defaults

The below tables describes the current attribute defaults for each implementation.

#### Search Base defaults

The following set defaults for the `additional_users_dn` and `additional_groups_dn` values.

| Implementation |   Users   |  Groups   |
|:--------------:|:---------:|:---------:|
|     lldap      | OU=people | OU=groups |

#### Attribute defaults

This table describes the attribute defaults for each implementation. i.e. the username_attribute is described by the
Username column.

| Implementation  |    Username    | Display Name | Mail | Group Name |
|:---------------:|:--------------:|:------------:|:----:|:----------:|
|     custom      |      N/A       | displayName  | mail |     cn     |
| activedirectory | sAMAccountName | displayName  | mail |     cn     |
|   rfc2307bis    |      uid       | displayName  | mail |     cn     |
|     freeipa     |      uid       | displayName  | mail |     cn     |
|      lldap      |      uid       |      cn      | mail |     cn     |
|     glauth      |       cn       | description  | mail |     cn     |

#### Filter defaults

The filters are probably the most important part to get correct when setting up LDAP. You want to exclude accounts under
the following conditions:

- The account is disabled or locked:
  - The [Active Directory] implementation achieves this via the `(!(userAccountControl:1.2.840.113556.1.4.803:=2))` filter.
  - The [FreeIPA] implementation achieves this via the `(!(nsAccountLock=TRUE))` filter.
  - The [GLAuth] implementation achieves this via the `(!(accountStatus=inactive))` filter.
  - The following implementations have no suitable attribute for this as far as we're aware:
    - [RFC2307bis]
    - [lldap]
- Their password is expired:
  - The [Active Directory] implementation achieves this via the `(!(pwdLastSet=0))` filter.
  - The [FreeIPA] implementation achieves this via the `(krbPasswordExpiration>={date-time:generalized})` filter.
  - The following implementations have no suitable attribute for this as far as we're aware:
    - [RFC2307bis]
    - [GLAuth]
    - [lldap]
- Their account is expired:
  - The [Active Directory] implementation achieves this via the `(|(!(accountExpires=*))(accountExpires=0)(accountExpires>={date-time:microsoft-nt}))` filter.
  - The [FreeIPA] implementation achieves this via the `(|(!(krbPrincipalExpiration=*))(krbPrincipalExpiration>={date-time:generalized}))` filter.
  - The following implementations have no suitable attribute for this as far as we're aware:
    - [RFC2307bis]
    - [GLAuth]
    - [lldap]

| Implementation  |                                                                                                                       Users Filter                                                                                                                       |                                                               Groups Filter                                                               |
|:---------------:|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:-----------------------------------------------------------------------------------------------------------------------------------------:|
|     custom      |                                                                                                                           N/A                                                                                                                            |                                                                    N/A                                                                    |
| activedirectory | (&(&#124;({username_attribute}={input})({mail_attribute}={input}))(sAMAccountType=805306368)(!(userAccountControl:1.2.840.113556.1.4.803:=2))(!(pwdLastSet=0))(&#124;(!(accountExpires=*))(accountExpires=0)(accountExpires>={date-time:microsoft-nt}))) |                               (&(member={dn})(&#124;(sAMAccountType=268435456)(sAMAccountType=536870912)))                                |
|   rfc2307bis    |                                                         (&(&#124;({username_attribute}={input})({mail_attribute}={input}))(&#124;(objectClass=inetOrgPerson)(objectClass=organizationalPerson)))                                                         | (&(&#124;(member={dn})(uniqueMember={dn}))(&#124;(objectClass=groupOfNames)(objectClass=groupOfUniqueNames)(objectClass=groupOfMembers))) |
|     freeipa     |   (&(&#124;({username_attribute}={input})({mail_attribute}={input}))(objectClass=person)(!(nsAccountLock=TRUE))(krbPasswordExpiration>={date-time:generalized})(&#124;(!(krbPrincipalExpiration=*))(krbPrincipalExpiration>={date-time:generalized})))   |                                                (&(member={dn})(objectClass=groupOfNames))                                                 |
|      lldap      |                                                                                 (&(&#124;({username_attribute}={input})({mail_attribute}={input}))(objectClass=person))                                                                                  |                                                (&(member={dn})(objectClass=groupOfNames))                                                 |
|     glauth      |                                                                 (&(&#124;({username_attribute}={input})({mail_attribute}={input}))(objectClass=posixAccount)(!(accountStatus=inactive)))                                                                 |                                              (&(uniqueMember={dn})(objectClass=posixGroup))                                               |

##### Microsoft Active Directory sAMAccountType

| Account Type Value |               Description               |               Equivalent Filter                |
|:------------------:|:---------------------------------------:|:----------------------------------------------:|
|     268435456      | Global/Universal Security Group Objects |                      N/A                       |
|     536870912      |   Domain Local Security Group Objects   |                      N/A                       |
|     805306368      |          Normal User Accounts           | `(&(objectCategory=person)(objectClass=user))` |

*__References:__*
- Account Type Values: [Microsoft Learn](https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-samr/e742be45-665d-4576-b872-0bc99d1e1fbe).
- LDAP Syntax Filters: [Microsoft TechNet Wiki](https://social.technet.microsoft.com/wiki/contents/articles/5392.active-directory-ldap-syntax-filters.aspx)
