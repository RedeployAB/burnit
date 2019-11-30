# secretgen

> Service with API to generate random character strings

`secretgen` provides an API to generate random character strings of varying
lengths, with or without special characters.

## Configuration

There are three ways of configuring the service. Either use environment variables, provide a config file
or use defaults.

**Service configuration**

**Environment variables**

* `SECRET_GEN_PORT` - Port the service listens to. Defaults to `3002`

**Configuration file**

Pass `-config` with path when running service, like so:
```
./secretgen -config config.yaml
```

*Example `config.yaml`*

```yaml
port: 3002

```
