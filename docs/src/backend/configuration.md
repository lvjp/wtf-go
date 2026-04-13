# Configuration

Default path as defined in the source code :

```go
{{#include ../../../internal/pkg/cmd/util/flags.go:default_config_path}}
```

## Example

This is the actual configuration file used for the local deployment.

```yaml
{{#include ../../../deployments/local/config.yaml}}
```

## Syntax

`log.level`:
:   Level used for logging.  
    Valid values: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`

`log.format`:
:   Logging output format.  
    Valid values :
    - `json`: JSON formatted output
    - `console`: Shiny debugging colored output for console

`server.listen_address`:
: Address used to listen for incoming HTTP requests with [fiber/App.Listen][fiber/App.Listen].

[fiber/App.Listen]: https://pkg.go.dev/github.com/gofiber/fiber/v3#App.Listen
