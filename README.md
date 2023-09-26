# hny-btgen
Honeycomb Board Template Generator. This tool will generate a board template
from an existing board within Honeycomb.

## Building

This will build the binary: `hny-btgen`
```shell
make build
```

## Usage

```shell
hny-btgen --honeycomb-api-key <HONEYCOMB_API_KEY> --board <BOARD_ID> [options]
```

### Options

The following options can be specified on the command line or via environment
variables. The Honeycomb API Key option is required and must be specified on the
or as an environment variable.

| CLI option          | Environment Variable | Description                                                  | Default |
|---------------------|----------------------|--------------------------------------------------------------|---------|
| --honeycomb-api-key | HONEYCOMB_API_KEY    | Honeycomb API Key with permissions to update dataset columns | `nil`   |
| --board             |                      | Honeycomb Board Id to use                                    | `nil`   |
| --out               |                      | Output template to file                                      | `nil`   |
| --graphic           |                      | Graphic # to use for the board template                      | `1`     |
| --sequence-number   |                      | Sequence # to use for board template, needs to be unique     | `99999` |
| --version           |                      | Display version information                                  | `false` |
