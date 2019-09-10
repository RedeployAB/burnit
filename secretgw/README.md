# secretgw

> Service with API to act as a gateway to `secretgen` and `secretdb`

## Configuration

The following environment variables needs to be set:

**Service configuration**

* `SECRET_GW_SERVICE_PORT` - Port the service listens to. Defaults to `3000`
* `SECRET_GENERATOR_SERVICE_BASE_URL` - URL with port to `secretgen`. Defaults to `http://localhost:3002`
* `SECRET_DB_SERVICE_BASE_URL` - URL with port to `secretdb`. Defaults to `http://localhost:3001`
* `SECRET_DB_SERVICE_API_KEY` - API key/token to `secretdb` endpoints
