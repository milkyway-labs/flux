# Configuration

This library expects to load a `config.yaml` file in order to initialize the required
`Indexer`, `Database`, `Node`, and `Module` instances.
In this section, you'll learn about the configuration structure and how the library
processes it to create the component instances used for indexing a chain.  

A complete configuration example can be found [here](./config-example.yaml).

## Overview

The configuration file includes the following sections:

* [logging](#logging): Defines the log format and verbosity level.
* [monitoring](#monitoring): Defines the prometheus exporter configuration.
* [databases](#databases): Contains configurations for databases that can be used by an `Indexer`.
* [nodes](#nodes): Contains configurations for the nodes that can be used by an `Indexer`.
* [modules](#modules): Contains module-specific configurations.
* [indexers](#indexers): Defines the configuration used to create `Indexer` instances.

### Logging 

Below is an example of a valid `logging` configuration:

```yaml
logging:
  format: "text"
  level: "debug"
```

Fields:

* `format`: Specifies the log format. Valid values are `text` and `json`.
* `level`: Specifies the log verbosity. For supported levels, refer to the [zerolog documentation](https://github.com/rs/zerolog/blob/9dacc014f38d60f563c2ab18719aec11fc06765c/globals.go#L36).

### Monitoring

Below is an example of a valid `monitoring` configuration:

```yaml
monitoring:
  enabled: true
  port: 2112
```

Fields:

* `enabled`: Specifies if the prometheus exporter should be enabled. Defaults to `true`.
* `port`: Port on which the prometheus exporter will listen. Defaults to `2112`.

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

### Indexers

Indexers are defined as a list, where each element represents the configuration for an `Indexer` instance. 
Since this is a component for which we know all possible configuration fields, we enforce the following structure:

```yaml
indexers:
  - name: "test-indexer"
    node_id: "cosmos-node"
    database_id: "my-db"
    # List of modules used by this indexer
    modules: 
      - "example"
    # Number of workers used to process blocks
    workers: 1 
    # Size of the queue used to track blocks to fetch
    height_queue_size: 100
    # Interval for polling the node for newly produced blocks
    node_polling_interval: "1s"
    # Number of attempts to retry parsing a block
    max_attempts: 2
    # Delay before re-enqueuing a failed block
    time_before_retry: "5s"
    # Indexer-specific module configurations
    override_module_config:
      example:
        value1: "custom"
        value2: 43
```

The following fields are **required**, while the rest are **optional** and will fall back to default values if not specified:

* `name`: The name of the indexer.
* `node_id`: The ID of the node used to fetch blocks.
* `database_id`: The ID of the database where indexed data will be stored.
* `modules`: A list of modules that will be used by the indexer.

**Optional fields:**

* `workers`: Number of workers used to process blocks. Defaults to `1`.
* `height_queue_size`: Maximum number of blocks that can be queued for fetching. Defaults to `100`.
* `node_polling_interval`: How often to poll the node for new blocks. Defaults to `"1s"`.
* `max_attempts`: Number of retries for parsing a failed block. Defaults to `5`
* `time_before_retry`: Delay before re-enqueuing a failed block for parsing. Defaults to `"10s"`
* `override_module_config`: A map containing module configurations specific to this indexer. 
This can be used to override the default configurations defined in the `modules` section.
* `start_height`: Height from which the indexer will start fetching blocks. If undefined the indexer will start indexing from the current node height.
* `force_reparse_old_blocks`: If `start_height` is defined, this flag will force the indexer to reparse the blocks from the start height to the current node height.
* `disabled`: If `true`, the indexer will not be started.

