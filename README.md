# hny-btgen
Honeycomb Board Template Generator. This tool will generate a board template
from an existing board within Honeycomb. The Short Description for each
query in the board template, will come from the Caption of the query in the
board. 

The board template will be written to stdout or to a file if the
`--out` option is specified. 

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
| --variables         |                      | Path to file with Variable Column mappings                   | `nil`   |
| --graphic           |                      | Graphic # to use for the board template                      | `1`     |
| --sequence-number   |                      | Sequence # to use for board template, needs to be unique     | `99999` |
| --version           |                      | Display version information                                  | `false` |

## Variables

You can specify a file with variable column mappings. The file can be in 
JSON or YAML format. The file should contain a map of the column name to an
array of value providers. The following JSON and YAML snippets are equivalent
for the variables file.

**JSON Format**

```json
{
  "variables": [
    {
      "name": "metrics.cpu.usage",
      "providers": [
        {
          "kind": "ExactMatch",
          "value": "metrics.cpu.usage"
        },
        {
          "kind": "ExactMatch",
          "value": "metrics_cpu_usage"
        }
      ]
    },
    {
      "name": "metrics.memory.utilization",
      "providers": [
        {
          "kind": "ExactMatch",
          "value": "metrics.memory.utilization"
        },
        {
          "kind": "ExactMatch",
          "value": "metrics_memory_utilization"
        }
      ]
    }
  ]
}
```

**YAML Format**

```yaml
variables:
  - name: metrics.cpu.usage
    providers:
      - kind: ExactMatch
        value: metrics.cpu.usage
      - kind: ExactMatch
        value: metrics_cpu_usage
  - name: metrics.memory.utilization
    providers:
      - kind: ExactMatch
        value: metrics.memory.utilization
      - kind: ExactMatch
        value: metrics_memory_utilization
```

The provider kind must be: `ExactMatch` or `SchemaMapping`
