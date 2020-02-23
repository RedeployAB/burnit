# burnitdb

> Service with API to handling secret requests for database

## Configuration

There are two ways of configuring the service. Either use environment
variables or provide a config file. Not all are mandatory, most
have default values.

***Service configuration**

**Environment variables**

* `SECRET_DB_PORT` - Port the service listens to. Defaults to `3001`
* `SECRET_DB_API_KEY` - API key/token to access the service endpoints (**mandatory**)
* `SECRET_DB_PASSPHRASE` - Passphrase for the hashes of the secret passphrases (**mandatory**)

*Database configuration*

* `DB_HOST` - FQDN/IP address for the MongoDB host. Defaults to `localhost`
* `DB` - Database for the secrets
* `DB_USER` - Database user with read/write access
* `DB_PASSWORD` - Password for the database user
* `DB_SSL` - True/False. If true, use SSL for DB communication. Defaults to `false`

**Configuration file**

Pass `-config` with path when running service, like so:
```
./burnitdb -config config.yaml
```

*Example `config.yaml`*

```yaml
server:
  port: 3001
  dbApiKey: <db-api-key> # Mandatory
  passphrase: secretstring # Mandatory
database:
  address: localhost:27017
  database: secret_db
  username: dbuser
  password: dbpassword
  ssl: true
  uri: mongodb://localhost:27017 # Can be used instead of address, database, username, password and ssl.
```
