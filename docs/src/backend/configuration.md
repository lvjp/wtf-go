# Configuration

## Loading

The configution can be loaded from several places in this order :

1. Builtin default values
2. File `/etc/wtf-go/config.yaml` if exists. JSON ???
3. File `$HOME/.config/wtf-go/config.yaml` is exists
4. The `--config` cli flag if available
5. Environment variables

## Defaults

```yaml
server:
  listen_address: :8080
log:
  level: info
  format: json
```

## Environment variables

Configuration is also loaded with configuration variable with prefix `WTF_GO_`.

For example, if you want to set `log.level`, you need to use environment variable:
`WTF_GO_LOG.LEVEL`.

## Syntax

`log.level`:
:   Level used for logging.  
    Default value: `info`  
    Valid values: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`  

`log.format`:
:   Logging output format.  
    Default value: `json`  
    Valid values :
      - `json`: JSON formatted output
      - `console`: Shiny debugging colored output for console

`server.listen_address`:
:   Address used to listen for incoming HTTP requests with [fiber/App.Listen][fiber/App.Listen].  
    Default value: `:8080`

[fiber/App.Listen]: https://pkg.go.dev/github.com/gofiber/fiber/v3#App.Listen
