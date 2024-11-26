# burnit

> Application for secret sharing.

<img src="assets/burnit-logo.png" alt="logo" width="300">

`burnit` is a service for creating temporary secrets and sharing them. In addition to this
it can be used to generate new secrets.

## Contents

* [Features](#features)
* [Requirements](#requirements)
* [Configuration](#configuration)
  * [Environment variables](#environment-variables)
* [TODO](#todo)

## Features

**Secret sharing**

The secrets are stored encrypted with 256-bit AES-GCM and are deleted upon retreival.
Either the encryption key (passphrase) kan be provided upon creation of a secret, or generated by the application.

**Secret generation**

The secret generation functionality returns a random string with length
and complexity based on the incoming request. These secrets are not stored.


## Configuration

There are four ways of configuring the service. Either provide a config file, use environment variables, pass command line arguments or use defaults.

Order of precedence:

* Defaults
* File
* Environment variables
* Command line arguments

### Environment variables**

| Name | Description |
|------|-------------|


## TODO

- [ ] Add CSRF protection
- [ ] Update documentation
- [ ] Add deployment examples, templates and scripts
