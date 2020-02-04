# burnit

> Repository for APIs and services for secret generation and sharing

This is a collection of three services:

* [`secretgw`](/secretgw/README.md) - Serves as a gateway and API endpoint for the other services
* [`secretgen`](/secretgen/README.md) - Generates random strings (secrets)
* [`secretdb`](/secretdb/README.md) - Stores random strings (secrets)

## API

These endpoints are available:

* `/api/v0/generate`
* `/api/v0/secrets`

### Generating secrets

```
// NOTE: special characters are: (!?=()&%)

// Generating a secret with default length (16).
GET /api/v0/generate

// Generating a secret with specified length of 24.
GET /api/v0/generate?length=24

// Generating a secret with special characters.
GET /api/v0/generate?specialchars=true

// Generating a secret with specified length and special characters.
GET /api/v0/generate?length=24&specialchars=true
```

### Creating and fetching secrets

```
// Creating a secret with no passphrase, and default TTL (7 days).
POST /api/v0/secrets
Body: {"secret":"<value>"}

// Creating a secret with a passphrase.
POST /api/v0/secrets
Body: {"secret":"<value>","passphrase":"<value>"}

// Creating a secret with a specified TTL in minutes.
POST /api/v0/secrets
Body: {"secret":"<value>","ttl":<value>}

// Creating a secret with a passphrase and a specified TTL in minutes.
POST /api/v0/secrets
Body: {"secret":"<value>","passphrase":"<value>","ttl":<value>}
```

**TODO**
Add more tests.
Determine on how to delete expired. Internal job in service, or side car.
Add date requirement on find.
