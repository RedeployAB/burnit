# burnit

> Repository for APIs and services for secret generation and sharing

This is a collection of three services:

* [`burnitgw`](/burnitgw/README.md) - Serves as a gateway and API endpoint for the other services
* [`burnitgen`](/burnitgen/README.md) - Generates random strings (secrets)
* [`burnitdb`](/burnitdb/README.md) - Stores random strings (secrets)

## API

These endpoints are available:

* `/generate`
* `/secrets`

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
