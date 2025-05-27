# Configuration

This library expects to load a `config.yaml` file in order to initialize the required
`Indexer`, `Database`, `Node`, and `Module` instances.
In this section, you'll learn about the configuration structure and how the library
processes it to create the component instances used for indexing a chain.  

A complete configuration example can be found [here](./config-example.yaml).

## Overview

The configuration file includes the following sections:

* `logging`: Defines the log format and verbosity level.
* `databases`: Contains configurations for databases that can be used by an `Indexer`.
* `nodes`: Contains configurations for the nodes that can be used by an `Indexer`.
* `modules`: Contains module-specific configurations.
* `indexers`: Defines the configuration used to create `Indexer` instances.

### Logging Configuration

Below is an example of a valid `logging` configuration:

```yaml
logging:
  format: "text"
  level: "debug"
```

Fields:

* `format`: Specifies the log format. Valid values are `text` and `json`.
* `level`: Specifies the log verbosity. For supported levels, refer to the [zerolog documentation](https://github.com/rs/zerolog/blob/9dacc014f38d60f563c2ab18719aec11fc06765c/globals.go#L36).

### Databases

Database configurations are defined as a map, where each key represents a unique database ID.
This ID can be referenced in the `database_id` field of an `Indexer` to link it with a specific database.

To keep the configuration flexible, the only required field for each database entry is `type`.
This allows for different database drivers to be registered and initialized with their specific YAML configurations.

For example, a PostgreSQL configuration might look like:

```yaml
databases:
  my-db:
    type: "postgres"
    host: "db-host"
    port: 5432
    database: "indexer"
    username: "username"
    password: "password"
```

During initialization, the following YAML content will be passed to the database driver:

```yaml
type: "postgres"
host: "db-host"
port: 5432
database: "indexer"
username: "username"
password: "password"
```

### Nodes

Nodes configurations are defined as a map, where each key represents a unique node ID.
This ID can be referenced in the `node_id` field of an `Indexer` to link it with a specific node.

To keep the configuration flexible, the only required field for each node entry is `type`.
This allows for different nodes to be registered and initialized with their specific YAML configurations.

For example, a Cosmos-SDK based node configuration might look like:

```yaml
nodes:
  cosmos-node:
    type: "cosmos"
    rpc_url: "https://rpc.chain.zone"
```

During initialization, the following YAML content will be passed to the node:

```yaml
type: "cosmos"
rpc_url: "https://rpc.chain.zone"
```

### Modules

Module configurations are defined as a map, where each key represents a unique module name.
This name can be included in the `modules` list of an `Indexer` to specify that the indexer should use that module.

Since custom modules may require different configurations, no specific structure is enforced for each map entry.

Here is an example of a module with two configuration values, which will be passed to the module during its initialization:

```yaml
modules:
  example:
    value1: "a generic value"
    value2: 42
```

During initialization, the module will receive the following YAML content:

```yaml
value1: "a generic value"
value2: 42
```

