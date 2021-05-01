# burnitgen

> Service with API to generate random character strings

`burnitgen` provides an API to generate random character strings of varying
lengths, with or without special characters.

## Configuration

There are four ways of configuring the service. Either provide a config file, use environment variables, pass command line arguments or use defaults.

Order of precedence:

* Defaults
* File
* Environment variables
* Command line arguments

**Service configuration**

**Environment variables**

* `BURNITGEN_LISTEN_PORT` - Port the service listens to. Defaults to `3002`

**Configuration file**

Pass `-config` with path when running service, like so:
```
./burnitgen -config config.yaml
```

*Example `config.yaml`*

```yaml
port: 3002

```

**Command line arguments**

```shell
  -config string
        Path to configuration file
  -port string
        Port to listen on
```
