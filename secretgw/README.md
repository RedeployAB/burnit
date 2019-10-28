# secretgw

> Service with API to act as a gateway to `secretgen` and `secretdb`

## Configuration

The following environment variables needs to be set:

**Service configuration**

* `SECRET_GW_PORT` - Port the service listens to. Defaults to `3000`
* `SECRET_GEN_BASE_URL` - URL with port to `secretgen`. Defaults to `http://localhost:3002`
* `SECRET_GEN_PATH` - Path for service calls. Defaults to `/api/v0/generate`
* `SECRET_DB_BASE_URL` - URL with port to `secretdb`. Defaults to `http://localhost:3001`
* `SECRET_DB_PATH`- Path for service alls. Defaults to `/api/v0/secrets`
* `SECRET_DB_API_KEY` - API key/token to `secretdb` endpoints
