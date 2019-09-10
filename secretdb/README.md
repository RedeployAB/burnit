# secretdb

> Service with API to handling secret requests for database

## Configuration

The following environment variables needs to be set:

**Service configuration**

* `SECRET_DB_SERVICE_PORT` - Port the service listens to. Defaults to `3001`
* `SECRET_DB_SERVICE_API_KEY` - API key/token to access the service endpoints
* `SECRET_DB_PASSPHRASE` - Passphrase for the hashes of the secret passphrases

**Database configuration**

* `DB_HOST` - FQDN/IP address for the MongoDB host. Defaults to `localhost`
* `DB` - Database for the secrets
* `DB_USER` - Database user with read/write access
* `DB_PASSWORD` - Password for the database user
* `DB_SSL` - True/False. If true, use SSL for DB communication. Defaults to `false`
