# burnit

> Application for secret generation and sharing.

## Contents

* [Introduction](#introduction)
* [Configuration](#configuration)
* [API](#api)
* [Usage](#usage)
* [Local development and testing](#loca-development-and-testing)
  * [docker-compose](#docker-compose)
* [Deploy to Kubernetes](#deploy-to-kubernetes)
* [Deploy to Azure Container Instances](#deploy-to-azure-container-instances)


## Introduction

`burnit` is a service for creating temporary secrets and sharing them. In addition to this
it can be used to generate new secrets.


**Secret generation**

The secret generation functionality returns a random string with length
and complexity based on the incoming request. These secrets are not stored.

**Secret sharing**

The secrets are stored encrypted with 256-bit AES-GCM and are deleted upon retreival.
Either the encryption key (passphrase) kan be provided upon creation of a secret,
and then passed along with the request to retreive it, or an encryption key
set in the application will be used.

## Configuration

There are four ways of configuring the service. Either provide a config file, use environment variables, pass command line arguments or use defaults.

Order of precedence:

* Defaults
* File
* Environment variables
* Command line arguments

**Environment variables**

* `BURNIT_LISTEN_HOST`
* `BURNIT_LISTEN_PORT`
* `BURNIT_LISTEN_ADDRESS`
* `BURNIT_ENCRYPTION_KEY`

*Database configuration*

### Running MongoDB in memory

To run a MongoDB in memory (or rather from a mounted RAM disk) issue
the following to your container:

```
mongod --smallfiles --noprealloc --nojournal --dbpath <ramdisk mounted localtion>
```

This information was got from the following [answer](https://stackoverflow.com/questions/26572248/can-i-configure-mongodb-to-be-in-memory) at stackoverflow.

**Kubernetes deployment**

```yaml
...
...
- name: mongo
  image: mongo:4.0
  ports:
  - containerPort: 27017
  command: [ "mongod" ]
  args: [ "--smallfiles", "--noprealloc", "--nojournal", "--dbpath",  "/data/inmem" ]
  volumeMounts:
  - mountPath: /data/inmem
    name: inmem
...
...
volumes:
- name: inmem
  emptyDir: {}
```

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
GET /secret?specialCharacters=true

// Generating a secret with specified length and special characters.
GET /secret?length=24&specialCharacters=true
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
  "ttl":<value>,          // Duration value. Such as 5m (5 minutes), 1h (1 hour) and 1d (1 day).
}

// Creating a secret with a passphrase and a specified TTL in minutes.
POST /secrets
Body:
{
  "value":"<value>",      // String value containing secret.
  "passphrase":"<value>", // String value containing passphrase.
  "ttl":<value>           // Duration value. Such as 5m (5 minutes), 1h (1 hour) and 1d (1 day).
}
```

**Retreiving**
```
// Retreive a secret without custom passphrase.
GET /secrets/<id>

// Retreive a secret with custom passphrase.
// Provide Passphrase in request headers.
HEADER:
Passphrase: <passphrase>

GET /secrets/<id>
```

**Example responses**
```
// Creating:
HEADER:
Location: /secrets/<id>

// (Uncompressed JSON for example readability)
{
  "id":"<id>",                    // String value containing id.
}

// Retreiving:
// (Uncompressed JSON for example readability)
{
  "id":"<id>",                    // String value containing id.
  "value":"<value>",              // String value containing secret.
}
```

## Local development and testing

## Deploy to Kubernetes

At `deployments/kubernetes` the following manifests are located:

```sh
deployments/kubernetes
├── burnit
│   ├── deployment.yaml
│   └── service.yaml
```

Create the following (minimal configuration) as their respective files named `config.yaml`:

```yaml
# Config
```

Deploy:

```sh
kubectl create namespace burnit

kubectl create secret generic burnit-config \
  --from-file=/path/to/file \
  --namespace burnit


kubectl apply -f deployments/kubernetes/burnit/deployment.yaml -n burnit
kubectl apply -f deployments/kubernetes/burnit/service.yaml -n burnit
```

Or run the provided script:
```sh
cd deployments/kubernetes
./deploy.sh --burnitdb-config /path/to/file --burnitgw-config /path/to/file
```

