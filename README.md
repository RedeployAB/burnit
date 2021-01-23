# burnit

> Application for secret generation and sharing.

## Contents

* [Introduction](#introduction)
* [API](#api)
* [Usage](#usage)
* [Local development and testing](#loca-development-and-testing)
  * [docker-compose](#docker-compose)
* [Deploy to Kubernetes](#deploy-to-kubernetes)


## Introduction

`burnit` is a collection of three services:

* [`burnitgw`](/burnitgw/README.md) - Serves as a gateway and API endpoint for the other services
* [`burnitgen`](/burnitgen/README.md) - Generates random strings (secrets)
* [`burnitdb`](/burnitdb/README.md) - Stores random strings (secrets)


**Secret generation**

The secret generation functionality returns a random string with length
and complexity based on the incoming request. These secrets are not stored.

**Secret sharing**

The secrets are stored encrypted with 256-bit AES-GCM and are deleted upon retreival.
Either the encryption key (passphrase) kan be provided upon creation of a secret,
and then passed along with the request to retreive it, or an encryption key
set in the application will be used.

There are two database options supported for the application. Redis and MongoDB,
with Redis being the default database in use.

For further instructions regarding configuration, check the `README` in each
of the services.

## API

These endpoints are available:

* `/secret`
* `/secrets`

## Usage

### Generating secrets

```
// NOTE: special characters are: _-!?=()&%

// Generating a secret with default length (16).
GET /secret

// Generating a secret with specified length of 24.
GET /secret?length=24

// Generating a secret with special characters.
GET /secret?specialchars=true

// Generating a secret with specified length and special characters.
GET /secret?length=24&specialchars=true
```

**Example response**
```
{"value":"<value>"}
```

### Creating and retreiving secrets

**Creating**
```
// Creating a secret with no passphrase, and default TTL (7 days).
POST /secrets
Body: {"value":"<value>"}

// Creating a secret with a passphrase.
POST /secrets
Body:
{
  "value":"<value>",      // String value containing secret.
  "passphrase":"<value>"  // String value containing passphrase.
}

// Creating a secret with a specified TTL in minutes.
POST /secrets
Body:
{
  "value":"<value>",      // String value containing secret.
  "ttl":<value>,          // Numerical value containing TTL in minutes.
}

// Creating a secret with a passphrase and a specified TTL in minutes.
POST /secrets
Body:
{
  "value":"<value>",      // String value containing secret.
  "passphrase":"<value>", // String value containing passphrase.
  "ttl":<value>           // Numerical value containing TTL in minutes.
}
```

**Retreiving**
```
// Retreive a secret without custom passphrase.
GET /secrets/<secretId>

// Retreive a secret with custom passphrase.
// Provide Passphrase in request headers.
HEADER:
Passphrase: <passphrase>

GET /secrets/<secretId>
```

**Example responses**
```
// Creating:
HEADER:
Location: /secrets/<id>

// (Uncompressed JSON for example readability)
{
  "id":"<id>",                    // String value containing id.
  "createdAt":"<date-and-time>",  // String value containing date and time in RFC3339.
  "expiresAt":"<date-and-time>"   // String value containing date and time in RFC3339.
}

// Retreiving:
// (Uncompressed JSON for example readability)
{
  "id":"<id>",                    // String value containing id.
  "value":"<value>",              // String value containing secret.
  "createdAt":"<date-and-time>",  // String value containing date and time in RFC3339.
  "expiresAt":"<date-and-time>"   // String value containing date and time in RFC3339.
}
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
  BURNITGEN_ADDRESS=http://burnitgen:3002
  BURNITDB_ADDRESS=http://burnitdb:3001

db.env
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

## Deploy to Kubernetes

At `deployments/kubernetes` the following manifests are located:

```sh
deployments/kubernetes
├── burnitdb
│   ├── deployment.yaml
│   └── service.yaml
├── burnitgen
│   ├── deployment.yaml
│   └── service.yaml
└── burnitgw
    ├── deployment.yaml
    └── service.yaml
```

The `deployment.yaml` for `burnitdb` and `burnitgw` expects
files as secrets to setup on the cluster for their
respective configuration.

Create the following (minimal configuration) as their respective files named `config.yaml`:

```yaml
# burnitdb
server:
  security:
    apiKey: <api-key-for-incoming-service-requests> # Optional.
    encryption:
      key: <encryption-key-string>
database:
  ssl: false
```
(`ssl: false` since `burnitdb` and `redis`/`mongo` reside in the same pod)

```yaml
# burnitgw
server:
  generatorAddress: burnitgen:3002
  dbAddress: burnitdb:3001
  dbApiKey: <same-api-key-as-above> # Optional if not set in DB.
```

Deploying:

```sh
kubectl create namespace burnit

kubectl create secret generic burnitdb-config \
  --from-file=/path/to/file \
  --namespace burnit

kubectl create secret generic burnitgw-config \
  --from-file=/path/to/file \
  --namespace burnit

# burnitgen
kubectl apply -f deployments/kubernetes/burnitgen/deployment.yaml -n burnit
kubectl apply -f deployments/kubernetes/burnitgen/service.yaml -n burnit
# burnitdb
kubectl apply -f deployments/kubernetes/burnitdb/deployment.yaml -n burnit
kubectl apply -f deployments/kubernetes/burnitdb/service.yaml -n burnit
# burnitgw
kubectl apply -f deployments/kubernetes/burnitgw/deployment.yaml -n burnit
kubectl apply -f deployments/kubernetes/burnitgw/service.yaml -n burnit
```

Or run the provided script:
```sh
cd deployments/kubernetes
./deploy.sh --burnitdb-config /path/to/file --burnitgw-config /path/to/file
```
