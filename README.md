<img src="assets/burnit-logo.png" alt="logo" width="300">

> Application for secret sharing.

`burnit` is a service for creating temporary secrets and sharing them. In addition to this
it can be used to generate new secrets.

## Contents

* [Features](#features)
* [Requirements](#requirements)
  * [Supported databases](#supported-databases)
* [Configuration](#configuration)
  * [Configuration file](#configuration-file)
  * [Environment variables](#environment-variables)
  * [Command-line flags](#command-line-flags)
* [Usage](#usage)
  * [API](#api)
    * [Secrets](#secrets)
    * [Errors](#errors)
      * [Error codes](#error-codes)
* [Sessions](#sessions)
* [Rate limiting](#rate-limiting)
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

- PostgreSQL
- MSSQL
- SQLite
- MongoDB
- Redis

## Configuration

There are four ways of configuring the service. Either provide a config file, use environment variables, pass command line arguments or use defaults.

Order of precedence:

* Defaults
* File
* Environment variables
* Command line arguments


### Configuration file

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
    # Database driver.
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
  runtimeRender: null
  # UI services configuration.
  services:
    session:
      # Timeout for the internal session service.
      timeout: 5s
      database:
        # Session database driver.
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

**Server** 

| Name | Description |
|------|-------------|
| `BURNIT_LISTEN_HOST` | Host (IP) to listen on. Default: `0.0.0.0`. |
| `BURNIT_LISTEN_PORT` | Port to listen on. Default: `3000`. |
| `BURNIT_TLS_CERT_FILE` | Path to TLS certificate file. |
| `BURNIT_TLS_KEY_FILE` | Path to TLS key file. |
| `BURNIT_CORS_ORIGIN` | CORS origin. Only necessary if frontend is not served through the server. |
| `BURNIT_RATE_LIMITER_RATE` | The average number of requests per second. |
| `BURNIT_RATE_LIMITER_BURST` | The maximum burst of requests. |
| `BURNIT_RATE_LIMITER_TTL` | The time-to-live for rate limiter entries. |
| `BURNIT_RATE_LIMITER_CLEANUP_INTERVAL` | The interval at which to clean up stale rate limiter entires. |
| `BURNIT_BACKEND_ONLY` | Disable UI (frontend). Default: `false`. |


**Secrets**

| Name | Description |
|------|-------------|
| `BURNIT_SECRET_SERVICE_TIMEOUT` | Timeout for the internal secret service. Default: `10s`. |



**UI**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_SERVICE_TIMEOUT` | Timeout for the internal session service. Default: `5s`. |
| `BURNIT_RUNTIME_RENDER` | Enable runtime rendering of the UI. |
| `BURNIT_LOCAL_DEVELOPMENT` | Enable local development mode. |

**Database**

| Name | Description |
|------|-------------|
| `BURNIT_DATABASE_DRIVER` | Database driver. |
| `BURNIT_DATABASE_URI` | URI (DSN) for the database. |
| `BURNIT_DATABASE_ADDRESS` | Address (host and port) for the database. |
| `BURNIT_DATABASE` | Database name. |
| `BURNIT_DATABASE_USER` | Database username. |
| `BURNIT_DATABASE_PASSWORD` | Database password. |
| `BURNIT_DATABASE_TIMEOUT` | Timeout for database operations. Default: `10s`. |
| `BURNIT_DATABASE_CONNECT_TIMEOUT` | Connect timeout for the database. Default: `10s`. |


**Database (MongoDB)**

| Name | Description |
|------|-------------|
| `BURNIT_DATABASE_MONGO_ENABLE_TLS` | Enable TLS for MongoDB. Default: true. |

**Database (Postgres)**

| Name | Description |
|------|-------------|
| `BURNIT_DATABASE_POSTGRES_SSL_MODE` | SSL mode for PostgreSQL. Default: require. |

**Database (MSSQL)**

| Name | Description |
|------|-------------|
| `BURNIT_DATABASE_MSSQL_ENCRYPT` | Encrypt for MSSQL. Default: true. |

**Database (SQLite)**

| Name | Description |
|------|-------------|
| `BURNIT_DATABASE_SQLITE_FILE` | Path to the database file for SQLite. Default: burnit.db. |
| `BURNIT_DATABASE_SQLITE_IN_MEMORY` | Use an in-memory database for SQLite. Default: false. |

**Database (Redis)**

| Name | Description |
|------|-------------|
| `BURNIT_DATABASE_REDIS_DIAL_TIMEOUT` | Dial timeout for the Redis client. |
| `BURNIT_DATABASE_REDIS_MAX_RETRIES` | Maximum number of retries for the Redis client. |
| `BURNIT_DATABASE_REDIS_MIN_RETRY_BACKOFF` |  Minimum retry backoff for the Redis client. |
| `BURNIT_DATABASE_REDIS_MAX_RETRY_BACKOFF` | Maximum retry backoff for the Redis client. |
| `BURNIT_DATABASE_REDIS_ENABLE_TLS` | Enable TLS for the Redis client. Default: true. |


**Session database**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_DATABASE_DRIVER` | Session database driver. |
| `BURNIT_SESSION_DATABASE_URI` | URI (DSN) for the session database. |
| `BURNIT_SESSION_DATABASE_ADDRESS` | Address (host and port) for the session database. |
| `BURNIT_SESSION_DATABASE` | Session database name. |
| `BURNIT_SESSION_DATABASE_USER` | Session Database username. |
| `BURNIT_SESSION_DATABASE_PASSWORD` | Session database password. |
| `BURNIT_SESSION_DATABASE_TIMEOUT` | Timeout for session database operations. Default: `5s`. |
| `BURNIT_SESSION_DATABASE_CONNECT_TIMEOUT` | Connect timeout for the session database. Default: `10s`. |


**Session database (MongoDB)**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_DATABASE_MONGO_ENABLE_TLS` | Enable TLS for MongoDB. Default: true. |

**Session database (Postgres)**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_DATABASE_POSTGRES_SSL_MODE` | SSL mode for PostgreSQL. Default: require. |

**Session database (MSSQL)**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_DATABASE_MSSQL_ENCRYPT` | Encrypt for MSSQL. Default: true. |

**Session database (SQLite)**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_DATABASE_SQLITE_FILE` | Path to the database file for SQLite. Default: burnit.db. |
| `BURNIT_SESSION_DATABASE_SQLITE_IN_MEMORY` | Use an in-memory database for SQLite. Default: false. |

**Session database (Redis)**

| Name | Description |
|------|-------------|
| `BURNIT_SESSION_DATABASE_REDIS_DIAL_TIMEOUT` | Dial timeout for the Redis client. |
| `BURNIT_SESSION_DATABASE_REDIS_MAX_RETRIES` | Maximum number of retries for the Redis client. |
| `BURNIT_SESSION_DATABASE_REDIS_MIN_RETRY_BACKOFF` |  Minimum retry backoff for the Redis client. |
| `BURNIT_SESSION_DATABASE_REDIS_MAX_RETRY_BACKOFF` | Maximum retry backoff for the Redis client. |
| `BURNIT_SESSION_DATABASE_REDIS_ENABLE_TLS` | Enable TLS for the Redis client. Default: true. |

### Command-line flags

```sh
Usage of burnit:
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
  -secret-service-timeout duration
        Optional. Timeout for the internal secret service. Default: 10s.
  -session-service-timeout duration
        Optional. Timeout for the internal session service. Default: 5s.
  -runtime-render value
        Optional. Enable runtime rendering of the UI.
  -local-development value
        Optional. Enable local development mode.
  -database-driver string
        Optional. Database driver.
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
  -session-database-driver string
        Optional. Database driver.
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

## Usage

### API

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

## Rate limiting

## TODO

- [ ] Update documentation
- [ ] Add deployment examples, templates and scripts
- [ ] Transactions for MongoDB commands
