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
* [Deploy to Kubernetes](#deploy-to-kubernetes)


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
databasE:
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
