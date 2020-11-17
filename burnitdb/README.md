# burnitdb

> Service with API to handling secret requests for database

`burnitdb` is a service for handling database access for
storing and handling secrets.

It supports two different databases: `redis` (default) and `mongodb` and
two hashing methods for the passphrase for protected secrets: `bcrypt` (default) and `md5`.

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

* `BURNITDB_LISTEN_PORT` - Port the service listens to. Defaults to `3001`
* `BURNITDB_API_KEY` - API key/token to access the service endpoints. Leave empty to not require an API key/token (**optional**)
* `BURNITDB_ENCRYPTION_KEY` - Encryption key for the secrets in the database (**mandatory**)
* `BURNITDB_HASH_METHOD` - Hash method for passphrases. Defaults to `bcrypt`

*Database configuration*

* `DB_HOST` - FQDN/IP address for the MongoDB host. Defaults to `localhost`
* `DB` - Database for the secrets
* `DB_USER` - Database user with read/write access
* `DB_PASSWORD` - Password for the database user
* `DB_SSL` - True/False. If true, use SSL for DB communication. Defaults to `true`
* `DB_DRIVER` - `redis`/`mongo`. The database engine to use for the service. Defaults to `redis`

**Configuration file**

Pass `-config` with path when running service, like so:
```
./burnitdb -config config.yaml
```

*Example `config.yaml`*

```yaml
server:
  port: 3001
  apiKey: <db-api-key> # Mandatory
  security:
    encryption:
      key: secretstring # Mandatory
database:
  driver: redis|mongo
  address: localhost:6379|localhost:27017
  database: burnit_db
  username: dbuser
  password: dbpassword
  ssl: true # Set to false if burnitdb and redis/mongo is in the same pod if using Kubernetes.
  uri: localhost:6379|mongodb://localhost:27017 # Can be used instead of address, database, username, password and ssl.
```

**Command line arguments**

```sh
  -api-key string
        API key for database endpoints
  -config string
        Path to configuration file
  -db string
        Database name
  -db-address string
        Host name and port for database
  -db-password string
        Password for user for database connections
  -db-uri string
        URI for database connection
  -db-user string
        User for database connections
  -disable-db-ssl
        Disable SSL for database connections
  -driver string
        Database driver for storage of secrets: mongo|redis
  -encryption-key string
        Encryption key for secrets in database
  -hash-method string
        Hash method for passphrase protected secrets
  -port string
        Port to listen onq
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
