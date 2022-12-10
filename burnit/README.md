# burnit

> Service with API to handling secret requests

`burnit` is a service for handling database access for
storing and handling secrets.

It supports two different databases: `redis` (default) and `mongo`.

* [Configuration](#configuration)
  * [Running MongoDB in memory](#running-mongodb-in-memory)


## Configuration

There are four ways of configuring the service. Either provide a config file, use environment variables, pass command line arguments or use defaults.

Order of precedence:

* Defaults
* File
* Environment variables
* Command line arguments

***Service configuration**

**Environment variables**

* `BURNIT_LISTEN_HOST` - Host (IP address) the service listens to. Defaults to `0.0.0.0`.
* `BURNIT_LISTEN_PORT` - Port the service listens to. Defaults to `3000`.
* `BURNIT_ENCRYPTION_KEY` - Encryption key for the secrets in the database (**mandatory**).
* `BURNIT_TLS_CERTIFICATE` - Path to TLS certificate file (`.crt`). Defaults to empty.
* `BURNIT_TLS_KEY` - Path to TLS key file (`.key`). Defaults to empty.
* `BURNIT_CORS_ORIGIN` - Enable CORS and sets `Access-Control-Allow-Origin` to provided value.


*Database configuration*

* `DB_HOST` - FQDN/IP address for the MongoDB host. Defaults to `localhost`.
* `DB` - Database for the secrets.
* `DB_USER` - Database user with read/write access.
* `DB_PASSWORD` - Password for the database user.
* `DB_SSL` - True/False. If true, use SSL for DB communication. Defaults to `true`.
* `DB_DRIVER` - `redis`/`mongo`. The database engine to use for the service. Defaults to `redis`.
* `DB_CONNECTION_URI` - URI for database connection.
* `DB_DIRECT_CONNECT` - Enable direct connect (mongodb only).

Use either `DB_CONNECTION_URI` or: `DB_HOST`, `DB`, `DB_USER`, `DB_PASSWORD`, `DB_SSL`.

**Configuration file**

Pass `-config` with path when running service, like so:
```
./burnit -config config.yaml
```

*Example `config.yaml`*

```yaml
server:
  host: "0.0.0.0"
  port: 3000
  security:
    encryption:
      key: secretstring # Mandatory
    tls:
      certificate: path/to/cert
      key: path/to/key
    cors:
      origin: <domain>
      
database:
  driver: redis|mongo
  address: localhost:6379|localhost:27017
  database: burnit_db
  username: dbuser
  password: dbpassword
  ssl: true # Set to false if burnit and redis/mongo is in the same pod if using Kubernetes.
  uri: localhost:6379|mongodb://localhost:27017 # Can be used instead of address, database, username, password and ssl.
  directConnect: false # Set to true when using MongoDB and direct connect is required.
```

**Command line arguments**

```sh
  -api-key string
        API key for database endpoints
  -config string
        Path to configuration file
  -cors-origin string
        Enable CORS and set origin
  -db string
        Database name
  -db-address string
        Host name and port for database
  -db-direct-connect
        Enable direct connect (mongodb only)
  -db-password string
        Password for user for database connections
  -db-uri string
        URI for database connection
  -db-user string
        User for database connections
  -disable-db-ssl
        Disable SSL for database connections
  -driver string
        Database driver for storage of secrets: redis|mongo
  -encryption-key string
        Encryption key for secrets in database
  -port string
        Port to listen on
  -tls-certificate string
        Path to TLS certificate file
  -tls-key string
        Path to TLS key file
```
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
