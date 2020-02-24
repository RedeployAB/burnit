# burnitgw

> Service with API to act as a gateway to `burnitgen` and `burnitdb`

## Configuration

There are two ways of configuring the service. Either use environment
variables or provide a config file. Not all are mandatory, most
have default values.

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
