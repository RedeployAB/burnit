<img src="assets/burnit-logo.png" alt="logo" width="300">

> Application for secret sharing.

`burnit` is a service for creating temporary secrets and sharing them. In addition to this
it can be used to generate new secrets.

## Contents

* [Features](#features)
* [Requirements](#requirements)
  * [Supported databases](#supported-databases)
  * [Running in a container](#running-in-a-container)
* [Install](#install)
* [Build](#build)
* [Configuration](#configuration)
  * [Configuration file](#configuration-file)
  * [Environment variables](#environment-variables)
  * [Command-line flags](#command-line-flags)
  * [Database configuration](#database-configuration)
    * [Database driver configuration](#database-driver-configuration)
* [Usage](#usage)
  * [API](#api)
    * [Secrets](#secrets)
    * [Errors](#errors)
      * [Error codes](#error-codes)
* [Sessions](#sessions)
* [Rate limiting](#rate-limiting)
* [Development](#development)
* [TODO](#todo)

## Features

**Secret sharing**

The secrets are stored encrypted with 256-bit AES-GCM and are deleted upon retreival.
Either the encryption key (passphrase) kan be provided upon creation of a secret, or generated by the application.

**Secret generation**

The secret generation functionality returns a random string with length
and complexity based on the incoming request. These secrets are not stored.

## Requirements

### Supported databases

The supported databases for the application are:

- PostgreSQL
- MSSQL
- SQLite
- MongoDB
- Redis
- In-memory

The main application database and the session database does not have to be the same database or driver.

### Running in a container

If running as a container the recommended resources are (as a start, and depending on expected load):

- **CPU**: `100m`
- **Memory**: `64Mi`


## Install

## Build

Scripts are provided to build the UI (frontend) and to build the application (that embeds the UI assets).

**Build UI (frontend)**

```sh
export ESBUILD_SHA256=77dce3e5d160db73bb37a61d89b5b38c5de1f18fbf4cc1c9c284a65ae5abb526
export HTMX_SHA256=24d3f3df3046e54d3fc260f58dcdeb82c53c38ee38f482eac74a5b6179d08ca7
export TAILWINDCSS_SHA256=cb5fff9d398d0d4f21c37c2e99578ada43766fbc6807b5f929d835ebfd07526b

./scripts/build-ui.sh
```

**Note**: The hashes listed above have been calculated by getting the checksum for:
- `esbuild@v0.24.0` from its official [download location](https://esbuild.github.io/dl/v0.24.0).
- `htmx@v2.0.3` from its official [download location](https://github.com/bigskysoftware/htmx/releases/download/v2.0.3/htmx.esm.js) + a modification to the file that replaces a call to `eval`.
- `tailwindcss@v3.4.14` from its official [download location](https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.14/tailwindcss-linux-x64).

The resulting files will be placed in `internal/ui/static` (they are present in `.gitignore`). This due to that they are to be embedded into the application binary.

**Build application**

```sh
# Build binary only:
./scripts/build.sh --version <version> --os <linux|darwin> --arch <amd64|arm64>

# Build binary and container image:
./scripts/build.sh --version <version> --os <linux|darwin> --arch <amd64|arm64> --image
```

The resulting binary (or container image) can then be run without any copying of assets due to them being embedded into the binary. Deploy the binary or container image to a target hosting environment
and run it.

## Configuration

There are four ways of configuring the service. Either provide a config file, use environment variables, pass command line arguments or use defaults.

Order of precedence:

* Defaults
* File
* Environment variables
* Command line arguments


In most scenarios the only configuration that needs to be set is the connection details for the database (main application).

Sessions and rate limiting uses in-memory databases due to their short lived nature by default. To avoid to use databases for these in scenarios where the application might need to scale, sticky sessions are an alternative.

**Example with configuration file**

The following configuration is all that is needed for most scenarios (example with PostgreSQL):

```yaml
services:
  secret:
    database:
      uri: postgres://<user>:<password>@<host>:5432/burnit
```

**Example with environment variables**

The following configuration is all that is needed for most scenarios (example with PostgreSQL):

```sh
export BURNIT_DATABASE_URI=postgres://<user>:<password>@<host>:5432/burnit
```

**Example with command-line flags**

The following configuration is all that is needed for most scenarios (example with PostgreSQL):

```sh
./burnit --database-uri postgres://<user>:<password>@<host>:5432/burnit
```

### Configuration file

By default the application will look for the file `config.yaml` in the same location as the binary. If another location and/or filename is desired, used the command-line flag: `-config-file <path>`.

All the available configuration that can be done with a configuration file:

```yaml
# Server configuration.
server:
  # Host (IP) to listen on.
  host: 0.0.0.0
  # Port to listen on.
  port: 3000
  tls:
    # Path to TLS certificate file.
    certFile: ""
    # Path to TLS key file.
    keyFile: ""
  cors:
    # CORS origin.
    # Only necessary if frontend is not served through the server.
    origin: ""
  # Rate limiting is disabled by default.
  # Configure one of the settings to enable it.
  rateLimiter:
    # The average number of requests per second.
    # Default: 1.
    rate: 0
    # The maximum burst of requests.
    # Default: 3.
    burst: 0
    # The time-to-live for rate limiter entries.
    # Default: 5m.
    ttl: 0s
    # The interval at which to clean up stale rate limiter entires.
    # Default: 10s.
    cleanupInterval: 0s
  # Disable UI (frontend).
  backendOnly: false 
# Service/application and database configuration.
services:
  secret:
    # Timeout for the internal secret service.
    timeout: 10s
    database:
      # Database driver. This is normally evaluated by the other database
      # configuration options but needs to be set if using a non-standard
      # port (when using address) or sqlite without options.
      driver: ""
      # URI (DSN) for the database.
      uri: ""
      # Address (host and port) for the database.
      address: ""
      # Database name.
      database: ""
      # Database username.
      username: ""
      # Database password.
      password: ""
      # Timeout for database operations.
      timeout: 10s
      # Connect timeout for the database.
      connectTimeout: 10s
      mongo:
        # Enable TLS for MongoDB.
        # Default: true.
        enableTLS: null
      postgres:
        # SSL mode for PostgreSQL.
        # Default: require.
        sslMode: ""
      mssql:
        # Encrypt for MSSQL.
        # Default: true.
        encrypt: ""
      sqlite:
        # Path to the database file for SQLite.
        # Default: burnit.db.
        file: ""
        # Use an in-memory database for SQLite.
        # Default: false.
        inMemory: null
      redis:
        # Dial timeout for the Redis client.
        # Default: 5s,
        dialTimeout: 0s
        # Maximum number of retries for the Redis client.
        # Default: 3,
        maxRetries: 0
        # Minimum retry backoff for the Redis client.
        # Default: 8ms.
        minRetryBackoff: 0s
        # Maximum retry backoff for the Redis client.
        # Default: 512ms.
        maxRetryBackoff: 0s
        # Enable TLS for the Redis client.
        # Default: true.
        enableTLS: null
# UI configuration.
ui:
  runtimeParse: null
  # UI services configuration.
  services:
    session:
      # Timeout for the internal session service.
      timeout: 5s
      database:
        # Session database driver. This is normally evaluated by the other
        # database configuration options but needs to be set if using a
        # non-standard port (when using address) or sqlite without options.
        driver: ""
        # URI (DSN) for the session database.
        uri: ""
        # Address (host and port) for the session database.
        address: ""
        # Session database name.
        database: ""
        # Session Database username.
        username: ""
        # Session database password.
        password: ""
        # Timeout for session database operations.
        timeout: 5s
        # Connect timeout for the session database.
        connectTimeout: 10s
        mongo:
          # Enable TLS for MongoDB.
          # Default: true.
          enableTLS: null
        postgres:
          # SSL mode for PostgreSQL. Default: require.
          sslMode: ""
        mssql:
          # Encrypt for MSSQL. Default: true.
          encrypt: ""
        sqlite:
          # Path to the database file for SQLite.
          # Default: burnit.db.
          file: ""
          # Use an in-memory database for SQLite.
          # Default: false.
          inMemory: null
        redis:
          # Dial timeout for the Redis client.
          # Default: 5s.
          dialTimeout: 0s
          # Maximum number of retries for the Redis client.
          # Default: 3.
          maxRetries: 0
          # Minimum retry backoff for the Redis client.
          # Default: 8ms.
          minRetryBackoff: 0s
          # Maximum retry backoff for the Redis client.
          # Default: 512ms.
          maxRetryBackoff: 0s
          # Enable TLS for the Redis client.
          # Default: true.
          enableTLS: null
```

### Environment variables

All the available configuration that can be done with environment variables:

**Server configuration**

| Name | Description |
|------|-------------|
| `BURNIT_LISTEN_HOST` | Host (IP) to listen on. Default: `0.0.0.0`. |
| `BURNIT_LISTEN_PORT` | Port to listen on. Default: `3000`. |
| `BURNIT_TLS_CERT_FILE` | Path to TLS certificate file. |
| `BURNIT_TLS_KEY_FILE` | Path to TLS key file. |
| `BURNIT_CORS_ORIGIN` | CORS origin. Only necessary if frontend is not served through the server. |
| `BURNIT_RATE_LIMITER` | Enable rate limiter with default values. Default: `false`. |
| `BURNIT_RATE_LIMITER_RATE` | The average number of requests per second. |
| `BURNIT_RATE_LIMITER_BURST` | The maximum burst of requests. |
| `BURNIT_RATE_LIMITER_TTL` | The time-to-live for rate limiter entries. |
| `BURNIT_RATE_LIMITER_CLEANUP_INTERVAL` | The interval at which to clean up stale rate limiter entires. |
| `BURNIT_BACKEND_ONLY` | Disable UI (frontend). Default: `false`. |


**Secrets configuration**

| Name | Description |
|------|-------------|
| `BURNIT_SECRET_SERVICE_TIMEOUT` | Timeout for the internal secret service. Default: `10s`. |


**Database configuration**

| Name | Description |
|------|-------------|
| `BURNIT_DATABASE_DRIVER` | Database driver. This is normally evaluated by the other database configuration options but needs to be set if using a non-standard port (when using address) or sqlite without options. |
| `BURNIT_DATABASE_URI` | URI (DSN) for the database. |
| `BURNIT_DATABASE_ADDRESS` | Address (host and port) for the database. |
| `BURNIT_DATABASE` | Database name. |
| `BURNIT_DATABASE_USER` | Database username. |
| `BURNIT_DATABASE_PASSWORD` | Database password. |
| `BURNIT_DATABASE_TIMEOUT` | Timeout for database operations. Default: `10s`. |
| `BURNIT_DATABASE_CONNECT_TIMEOUT` | Connect timeout for the database. Default: `10s`. |


**Database (MongoDB) configuration**

| Name | Description |
|------|-------------|
| `BURNIT_DATABASE_MONGO_ENABLE_TLS` | Enable TLS for MongoDB. Default: true. |

**Database (Postgres) configuration**

| Name | Description |
|------|-------------|
| `BURNIT_DATABASE_POSTGRES_SSL_MODE` | SSL mode for PostgreSQL. Default: require. |

**Database (MSSQL) configuration**

| Name | Description |
|------|-------------|
| `BURNIT_DATABASE_MSSQL_ENCRYPT` | Encrypt for MSSQL. Default: true. |

**Database (SQLite) configuration**

| Name | Description |
|------|-------------|
| `BURNIT_DATABASE_SQLITE_FILE` | Path to the database file for SQLite. Default: burnit.db. |
| `BURNIT_DATABASE_SQLITE_IN_MEMORY` | Use an in-memory database for SQLite. Default: false. |

**Database (Redis) configuration**

| Name | Description |
|------|-------------|
| `BURNIT_DATABASE_REDIS_DIAL_TIMEOUT` | Dial timeout for the Redis client. |
| `BURNIT_DATABASE_REDIS_MAX_RETRIES` | Maximum number of retries for the Redis client. |
| `BURNIT_DATABASE_REDIS_MIN_RETRY_BACKOFF` |  Minimum retry backoff for the Redis client. |
| `BURNIT_DATABASE_REDIS_MAX_RETRY_BACKOFF` | Maximum retry backoff for the Redis client. |
| `BURNIT_DATABASE_REDIS_ENABLE_TLS` | Enable TLS for the Redis client. Default: true. |


**UI configuration**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_SERVICE_TIMEOUT` | Timeout for the internal session service. Default: `5s`. |
| `BURNIT_RUNTIME_PARSE` | Enable runtime parsing of the UI templates. |


**Session database configuration**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_DATABASE_DRIVER` | Session database driver. This is normally evaluated by the other database configuration options but needs to be set if using a non-standard port (when using address) or sqlite without options. |
| `BURNIT_SESSION_DATABASE_URI` | URI (DSN) for the session database. |
| `BURNIT_SESSION_DATABASE_ADDRESS` | Address (host and port) for the session database. |
| `BURNIT_SESSION_DATABASE` | Session database name. |
| `BURNIT_SESSION_DATABASE_USER` | Session Database username. |
| `BURNIT_SESSION_DATABASE_PASSWORD` | Session database password. |
| `BURNIT_SESSION_DATABASE_TIMEOUT` | Timeout for session database operations. Default: `5s`. |
| `BURNIT_SESSION_DATABASE_CONNECT_TIMEOUT` | Connect timeout for the session database. Default: `10s`. |


**Session database (MongoDB) configuration**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_DATABASE_MONGO_ENABLE_TLS` | Enable TLS for MongoDB. Default: true. |

**Session database (Postgres) configuration**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_DATABASE_POSTGRES_SSL_MODE` | SSL mode for PostgreSQL. Default: require. |

**Session database (MSSQL) configuration**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_DATABASE_MSSQL_ENCRYPT` | Encrypt for MSSQL. Default: true. |

**Session database (SQLite) configuration**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_DATABASE_SQLITE_FILE` | Path to the database file for SQLite. Default: burnit.db. |
| `BURNIT_SESSION_DATABASE_SQLITE_IN_MEMORY` | Use an in-memory database for SQLite. Default: false. |

**Session database (Redis) configuration**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_DATABASE_REDIS_DIAL_TIMEOUT` | Dial timeout for the Redis client. |
| `BURNIT_SESSION_DATABASE_REDIS_MAX_RETRIES` | Maximum number of retries for the Redis client. |
| `BURNIT_SESSION_DATABASE_REDIS_MIN_RETRY_BACKOFF` |  Minimum retry backoff for the Redis client. |
| `BURNIT_SESSION_DATABASE_REDIS_MAX_RETRY_BACKOFF` | Maximum retry backoff for the Redis client. |
| `BURNIT_SESSION_DATABASE_REDIS_ENABLE_TLS` | Enable TLS for the Redis client. Default: true. |

### Command-line flags

All the available configuration that can be done with environment variables:

```sh
Command-line configuration for burnit:

  # Server configuration.
  -config-path string
        Optional. Path to a configuration file. Defaults to: config.yaml.
  -host string
        Optional. Host (IP) to listen on. Default: 0.0.0.0.
  -port int
        Optional. Port to listen on. Default: 3000.
  -tls-cert-file string
        Optional. Path to TLS certificate file.
  -tls-key-file string
        Optional. Path to TLS key file.
  -cors-origin string
        Optional. CORS origin for application. Only necessary if UI is not served through this application.
  -rate-limiter-burst int
        Optional. The maximum burst of requests.
  -rate-limiter-cleanup-interval duration
        Optional. The interval at which to clean up stale rate limiter entires.
  -rate-limiter-rate float
        Optional. The average number of requests per second.
  -rate-limiter-ttl duration
        Optional. The time-to-live for rate limiter entries.
  # Secrets configuration.
  -secret-service-timeout duration
        Optional. Timeout for the internal secret service. Default: 10s.
  -database-driver string
        Optional. Database driver. This is normally evaluated by the other database configuration options but needs to be set if using a non-standard port (when using address) or sqlite without options.
  -database-uri string
        Optional. URI (DSN) for the database.
  -database-address string
        Optional. Address (host and port) for the database.
  -database string
        Optional. Database name.
  -database-user string
        Optional. Database username.
  -database-password string
        Optional. Database password.
  -database-timeout duration
        Optional. Timeout for database operations. Default: 10s.
  -database-connect-timeout duration
        Optional. Connect timeout for the database. Default: 10s.
  -database-mongo-enable-tls value
        Optional. Enable TLS for MongoDB. Default: true.
  -database-postgres-ssl-mode string
        Optional. SSL mode for PostgreSQL. Default: require.
  -database-mssql-encrypt string
        Optional. Encrypt for MSSQL. Default: true.
  -database-sqlite-file string
        Optional. Path to the database file for SQLite. Default: burnit.db.
  -database-sqlite-in-memory value
        Optional. Use an in-memory database for SQLite. Default: false.
  -database-redis-dial-timeout duration
        Optional. Dial timeout for the Redis client.
  -database-redis-enable-tls value
        Optional. Enable TLS for the Redis client. Default: true.
  -database-redis-max-retries int
        Optional. Maximum number of retries for the Redis client.
  -database-redis-max-retry-backoff duration
        Optional. Maximum retry backoff for the Redis client.
  -database-redis-min-retry-backoff duration
        Optional. Minimum retry backoff for the Redis client.
  # UI configuration.
  -session-service-timeout duration
        Optional. Timeout for the internal session service. Default: 5s.
  -runtime-parse value
        Optional. Enable runtime parsing of the UI.
  -session-database-driver string
        Optional. Database driver. This is normally evaluated by the other database configuration options but needs to be set if using a non-standard port (when using address) or sqlite without options.
  -session-database-uri string
        Optional. URI for the session database.
  -session-database-address string
        Optional. Address for the session database.
  -session-database string
        Optional. Session database name.
  -session-database-user string
        Optional. Session database username.
  -session-database-password string
        Optional. Session database password.
  -session-database-timeout duration
        Optional. Timeout for session database operations. Default: 10s.
  -session-database-connect-timeout duration
        Optional. Connect timeout for the session database. Default: 10s.
  -session-database-mongo-enable-tls value
        Optional. Enable TLS for MongoDB. Default: true.
  -session-database-postgres-ssl-mode string
        Optional. SSL mode for PostgreSQL. Default: require.
  -session-database-mssql-encrypt string
        Optional. Encrypt for MSSQL. Default: true.
  -session-database-sqlite-file string
        Optional. Path to the database file for SQLite. Default: burnit.db.
  -session-database-sqlite-in-memory value
        Optional. Use an in-memory database for SQLite. Default: false.
  -session-database-redis-dial-timeout duration
        Optional. Dial timeout for the Redis client.
  -session-database-redis-enable-tls value
        Optional. Enable TLS for the Redis client. Default: true.
  -session-database-redis-max-retries int
        Optional. Maximum number of retries for the Redis client.
  -session-database-redis-max-retry-backoff duration
        Optional. Maximum retry backoff for the Redis client.
  -session-database-redis-min-retry-backoff duration
        Optional. Minimum retry backoff for the Redis client.
```

### Database configuration

The application supports various database drivers as mentioned in the [requirements](#requirements) section. The main database (containing secrets) and the database for handling sessions can be handled separately and does not need to be the same database or driver.

If not database configuration is set for the application it will default to using a built-in in-memory database. This will not persist secrets between restarts and is not recommended unless this is desired. Using the in-memory database will log a warning. Expired secrets are cleaned up from the database.

As with the main application database, if no database configuration is set for the session database it will default to using a built-in in-memory database. This is considered a normal configuration due to the lifetime cycle and nature of the sessions in this application and will
not log a warning. Expired sessions are cleaned up from the database.

#### Database driver configuration

In most circumstances the database driver does not need to be configured as it is evaluated by the database configuration.

The situations where the driver needs to be configured are:

* Using a non-standard port when using an address for database configuration.
* Using SQLite without specifying either a path to a database file, or specifying to use SQLite with in-memory mode.

The supported values for the database driver are:

* `mongodb`
* `postgres`
* `sqlserver`
* `sqlite`
* `redis`
* `inmem`


## Usage

### API

### Secret

#### Generate secret

```http
GET /secret
```

```http
GET /secret?length=32&specialCharacters=true
```

##### Headers

| Name | Required | Description |
|------|----------|-------------|
| `Accept` | **False** | Supported values: `application/json` and `plain/text`|

##### URI parameters

| Name | In | Required | Type | Description |
|------|----|----------|------|-------------|
| `length` | Query | **False** | *number* | Amount of characters in the secret. Default: `16`. Alias `l` can be used. |
| `specialCharacters` | Query| **False** | *boolean* | Use special characters or not. Default `false`. Alias `sc` can be used. |

**Note**: If not len

##### Response

```http
200 Status Ok
```

**`application/json`**

```json
{
  "value": "secret"
}
```

**`plain/text`**

```http
secret
```

### Secrets

#### Get secret

```http
GET /secrets/{id}
```

##### Headers

| Name | Required | Description |
|------|----------|-------------|
| `Passphrase` | **True** | Passphrase for the secret. |

##### URI parameters

| Name | In | Required | Type | Description |
|------|----|----------|------|-------------|
| `id` | Path | **True** | *string* | The ID of the secret to retrieve. |

##### Response

```json
{
  "value": "secret"
}
```

#### Create secret

##### Request body

```json
{
  "value": "secret",
  "passphrase": "passphrase",
  "ttl": "1h",
  "expiresAt": "2025-01-24T18:09:55+01:00"
}
```

| Name | Required | Type | Description |
| ---- | -------- | ---- | ----------- |
| `value` | **True** | *string* | Secret value. |
| `passphrase` | **False** | *string* | Passphrase for the secret. <sup>*1)</sup> |
| `ttl` | **False** | *string* | A time duration. Example: `1h`. <sup>*2)</sup><sup>*3)</sup><sup>*4)</sup> |
| `expiresAt` | **False** | *Date* | Date in RFC3399 (ISO 8601). Takes precedence over `ttl`. See example body. <sup>*3)</sup><sup>*4)</sup> |

**Note**

<sup>*1) A passphrase will be generated if non is provided.<br/>
<sup>*2) A duration according to the Go duration format. Example: `1m`, `1h` and so on. The highest unit is `h`. For 3 days the value should be `72h`. Can be used with additional units like so: `1h10m10s` which is 1 hour, 10 minutes and 10 seconds.</sup><br/>
<sup>*3) If neither `ttl` or `expiresAt` is provided a default expiration time of `1h` will be set.</sup><br/>
<sup>*4)Minumum expiration time is `1m` (1 minute) and maximum expiration time is `168h` (7 days).</sup>

##### Response

```http
201 Status Created
```

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "passphrase": "passphrase",
  "ttl": "1h0m0s",
  "expiresAt": "2025-01-24T18:09:55+01:00"
}
```

### Errors

Error responses have the following structure:

```json
{
  "statusCode": 400,
  "code": "ErrorCode",
  "error": "an error occured",
  "requestId": "00000000-0000-0000-0000-000000000000"
}
```

| Name | Type | Description |
|------|------|-------------|
| `statusCode` | *number* | HTTP status code of the error. |
| `code` | *string* | [Error code](#error-codes) of the error. |
| `error` | *string* | The error text/information/message. |
| `requestId` | *string* | The request ID (UUID) for the request triggering the error. |


#### Error codes

| Error code | HTTP status code | Description |
|------------|------------------|-------------|
| `EmptyRequest` | `400` | Request body for creating a secret is empty. |
| `InvalidRequest` | `400` | Request body for creating a secret is invalid. |
| `MalformedRequest` | `400` | Request body for creating a secret is malformed. |
| `PassphraseNotBase64` | `400` | Passphrase for a secret is not Base 64 encoded. |
| `InvalidExpirationTime` | `400` | Expiration time for secret is invalid. |
| `ValueInvalid` | `400` | Value for secret contains invalid characters, or has an invalid format. |
| `ValueTooManyCharacters` | `400` | Value for secret contains too many characters. |
| `PassphraseInvalid` | `400` | Passphrase for secret contains invalid characters, or has an invalid format. |
| `PassphraseTooFewCharacters` | `400` | Passphrase has too few characters. |
| `PassphraseTooManyCharacters` | `400` | Passphrase has too many characters. |
| `InvalidBase64` | `400` | `400` | Invalid Base 64 encoded string provided. |
| `ErrPassphraseRequired` | `401` | Passphrase required. |
| `InvalidPassphrase` | `401` | Passphrase for secret is invalid. |
| `SecretNotFound` | `404` | Secret not found. Either secret does not exist, or has been read. |

## Sessions

The application handle sessions with CSRF tokens to increase security when creating and retrieving secrets. Sessions are short-lived with a lifetime of 15 minutes. The application clears out expired sessions from the database every minute which frees up memory.

It is also possible to store sessions in a database. See more at the sections [Database configuration](#database-configuration), [Configuration file](#configuration-file), [Environment variables](#environment-variables) and [Command-line flags](#command-line-flags).


## Rate limiting

A simple rate limiting mechanism is built-in into the application. It handles rate limiting on a per IP basis and store the data in an in-memory database. The rate limiting model is according to a token bucket algorithm that allows for requests to be made as long as there are tokens in the bucket.

If a request is made it will refill with *n* tokens per second, with an allowed burst of *n*.
A rate limit for an IP address have a default time-to-live of 5 minutes and are cleared out periodically (default every 10 seconds).

**Example of rate limiting**

Rate is configured to `1` and burst is configured to `3` it will allow for an average of 1 request per second with a maximum burst of 3 in consecutive requests.

### Enable rate limiting

To enable rate limiting either set one or more options, or set the environment `BURNIT_RATE_LIMITER=true`, use the command-line flag `-rate-limiter=true` or enable it in the config file:

```yaml
server:
  rateLimiter:
    enabled: true
```


The options that are not configured will have the following default values:

* Rate: `1`
* Burst: `3`
* TTL: `5m`
* Cleanup interval `10s`

If more advanced rate limiting is required, do not enable rate limiting and configure an external rate limiter.

## Development

To develop the application the following tools are needed:

* `tailwindcss` - Build the CSS file(s).
* `esbuild` - Bundle and minify CSS and JavaScript files.

**Note**: Only required if developing the UI/frontend.

To actively build the CSS files while making modifications to the HTML templates run `tailwindcss` with watching:

```sh
cd internal/ui
tailwindcss -i ./static/css/tailwind.css -o ./static/css/main.css --watch
```

Set `BURNIT_RUNTIME_PARSE=true`, use the command-line flag `--runtime-parse=true` or enable it in the config file:

```yaml
ui:
  runtimeParse: true
```

This will make sure the application parses the HTML template every call, thus making it possible to see changes to HTML templates, JavaScript and CSS at every save.

## TODO

- [ ] Add deployment examples, templates and scripts
