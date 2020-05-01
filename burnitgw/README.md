# burnitgw

> Service with API to act as a gateway to `burnitgen` and `burnitdb`

## Configuration

There are four ways of configuring the service. Either provide a config file, use environment variables, pass command line arguments or use defaults.

Order of precedence:

* Defaults
* File
* Environment variables
* Command line arguments

**Service configuration**

**Environment variables**

* `BURNITGW_LISTEN_PORT` - Port the service listens to. Defaults to `3000`
* `BURNITGEN_BASE_URL` - URL with port to `burnitgen`. Defaults to `http://localhost:3002`
* `BURNITGEN_PATH` - Path for service calls. Defaults to `/api/generate`
* `BURNITDB_BASE_URL` - URL with port to `burnitdb`. Defaults to `http://localhost:3001`
* `BURNITDB_PATH`- Path for service alls. Defaults to `/api/secrets`
* `BURNITDB_API_KEY` - API key/token to `burnitdb` endpoints

**Configuration file**

Pass `-config` with path when running service, like so:
```
./burnitgw -config config.yaml
```

*Example `config.yaml`*

```yaml
server:
  port: 3000
  generatorBaseUrl: "http://localhost:3002"
  generatorServicePath: "/api/generate"
  dbBaseUrl: "http://localhost:3001"
  dbServicePath: "/api/secrets"
  dbApiKey: "<DB-API-KEY>"
```

**Command line arguments**

```
  -config string
        Path to configuration file
  -db-api-key string
        API Key to DB service
  -db-base-url string
        Base URL to DB service (burnitdb)
  -db-service-path string
        Path to DB service endpoint (burnitdb)
  -generator-base-url string
        Base URL to generator service (burnitgen)
  -generator-service-path string
        Path to generator service endpoint (burnitgen)
  -port string
        Port to listen on
```
