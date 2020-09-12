# burnit

> Repository for APIs and services for secret generation and sharing

This is a collection of three services:

* [`burnitgw`](/burnitgw/README.md) - Serves as a gateway and API endpoint for the other services
* [`burnitgen`](/burnitgen/README.md) - Generates random strings (secrets)
* [`burnitdb`](/burnitdb/README.md) - Stores random strings (secrets)

## Contents

* [API](#api)
* [Usage](#usage)
* [Local development and testing](#loca-development-and-testing)
  * [docker-compose](#docker-compose)


## API

These endpoints are available:

* `/generate`
* `/secrets`

## Usage

### Generating secrets

```
// NOTE: special characters are: (!?=()&%)

// Generating a secret with default length (16).
GET /generate

// Generating a secret with specified length of 24.
GET /generate?length=24

// Generating a secret with special characters.
GET /generate?specialchars=true

// Generating a secret with specified length and special characters.
GET /generate?length=24&specialchars=true
```

### Creating and fetching secrets

```
// Creating a secret with no passphrase, and default TTL (7 days).
POST /secrets
Body: {"secret":"<value>"}

// Creating a secret with a passphrase.
POST /secrets
Body: {"secret":"<value>","passphrase":"<value>"}

// Creating a secret with a specified TTL in minutes.
POST /secrets
Body: {"secret":"<value>","ttl":<value>}

// Creating a secret with a passphrase and a specified TTL in minutes.
POST /secrets
Body: {"secret":"<value>","passphrase":"<value>","ttl":<value>}
```

## Local development and testing

### docker-compose

In `deployments/compose` a `docker-compose.yml` is located.
To get started define the following files in the same directory:

```
.env:
  BURNITGW_VERSION=<tag>
  BURNITDB_VERSION=<tag>
  BURNITGEN_VERSION=<tag>
  DB_IMAGE=redis|mongo
  DB_PORT=6379|27017

gw.env:
  BURNITGEN_BASE_URL=http://burnitgen:3002
  BURNITDB_BASE_URL=http://burnitdb:3001
  BURNITDB_API_KEY=<string>

db.env
  BURNITDB_API_KEY=<string>
  BURNITDB_ENCRYPTION_KEY=<string>
  DB_HOST=redis|mongo
  DB_SSL=false
  DB_DRIVER=redis|mongo
```

**Details**

* `BURNITDB_API_KEY` in `db.env` and `gw.env` should be the same
string
* Configure for redis:
  * `.env`
    * `DB_IMAGE=redis`
    * `DB_PORT=6379`
  * `db.env`
    * `DB_HOST=redis`
    * `DB_DRIVER=redis`
* Configure for mongodb:
  * `.env`
    * `DB_IMAGE=mongo`
    * `DB_PORT=27017`
  * `db.env`
    * `DB_HOST=mongo`
    * `DB_DRIVER=mongo`
