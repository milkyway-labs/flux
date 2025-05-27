# Support new chains

This document explains how to support a new blockchain.
The library is designed to isolate blockchain-specific details in the `Node`, `Block`, and `Tx` interfaces inside the `types` module.
To support a new chain, a developer must define custom types that implement these interfaces.

---

## Create a New Tx Type

A `Tx` represents an action performed by a user that has been included in a `Block`.
This object should contain transaction-related information that may be useful during the indexing process.

To define a custom `Tx` type, implement the `Tx` interface provided by the `types` module. Here's its definition:

```go
// Tx represents a generic transaction that has been included in a Block.
type Tx interface {
	// GetHash returns the transaction hash.
	GetHash() string
	// IsSuccessful returns true if the transaction was executed without errors.
	IsSuccessful() bool
}
```

Once your custom `Tx` type is defined, you can proceed to create your `Block` type.

For an example of a custom `Tx` implementation, refer to the [Cosmos SDK-based implementation](../cosmos/types/chain.go).

---

## Create a New Block Type

Blockchains organize their data in a sequence of blocks. This structure should represent a block that has been included in the chain and should expose the transactions and any other chain-specific data that may be relevant for indexing.

To define a custom `Block` type, implement the following interface:

```go
// Block represents a generic block produced by a blockchain.
type Block interface {
	// GetChainID returns the ID of the blockchain that produced this block.
	GetChainID() string
	// GetHeight returns the height at which this block was produced.
	GetHeight() Height
	// GetTimeStamp returns the time at which this block was produced.
	GetTimeStamp() time.Time
	// GetTxs returns the transactions included in this block.
	GetTxs() []Tx
}
```

**Note:** Design your `Block` type to be efficiently retrievable by the `Node`, which is responsible for fetching and returning this structure.

By implementing the `Block` interface, the `Indexer` can fetch your blocks and send them to the modules for processing.

For an example, see the [Cosmos SDK-based block implementation](../cosmos/types/chain.go).

---

## Create a New Node

The final component to implement is the `Node`. It is responsible for fetching `Block`s from the chain so that modules can extract relevant data.

To define a custom `Node`, implement the following interface:

```go
// Node represents a generic blockchain node that can be queried to obtain Blocks.
type Node interface {
	// GetChainID returns the ID of the blockchain being queried.
	GetChainID() string
	// GetBlock fetches the block produced at the given height.
	GetBlock(context context.Context, height types.Height) (types.Block, error)
	// GetLowestHeight returns the lowest queryable block height from the node.
	GetLowestHeight(context context.Context) (types.Height, error)
	// GetCurrentHeight returns the current height of the node.
	GetCurrentHeight(context context.Context) (types.Height, error)
}
```

For an example, refer to the [Cosmos SDK-based node implementation](../cosmos/node/rpc/node.go).

**Note:** The `GetChainID()` function does not accept a `context.Context` parameter because it is not intended to perform a network request.
Each `Node` instance in this library is expected to index only a specific chain, 
so the chain ID can be cached during initialization and returned from memory or 
loaded from the node’s configuration.

If your `Node` requires configuration, define it using a Go `struct` that can be
parsed from a YAML file. 
This library expects all configurations to be defined in a `config.yaml` file.

For example:

```go
type Config struct {
    URL string `yaml:"url"`
}
```

During initialization, the library will call a builder function 
with the YAML configuration provided as a `[]byte`, which you can unmarshal into your custom `Config` struct.

---

## Register your Node type

To register your custom `Node` and allow the library to build an instance of it,
you must provide a `Builder` function with the following signature:

```go
// Builder represents a function used to build a Node instance.
// It receives the raw YAML configuration bytes and an ID that identifies the node.
type Builder func(ctx context.Context, id string, rawConfig []byte) (node.Node, error)
```

This allows the `NodesManager` to construct your node instance when an `Indexer` requires it.

To register your node, call the `NodesManager`'s `RegisterNode(type string, builder Builder)` function,
providing the node type identifier and the corresponding `Builder` function.
The `type` is used in the configuration file to determine which builder should be used to initialize the given `Node` type.

For example, if you’ve created an `EVMNode` for EVM-compatible chains and registered it with `type` set to `"evm"`,
the library will try to build an `EVMNode` whenever it finds an entry like the following in the [`nodes` section](./config_structure.md#nodes) of the configuration:

```yaml
type: "evm"
# other values that will be passed as configuration during build time
```

For a reference implementation, see the `Builder` function for the Cosmos SDK-based `Node` [here](../cosmos/node/rpc/builder.go).

