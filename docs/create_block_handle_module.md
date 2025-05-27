# Custom Module to Parse Data from a Block

This document explains how to create a custom indexing `Module` that extracts data from a `Block`.

---

## Implement the `BlockHandleModule` Interface

To indicate that your module can extract data from a `Block`, your struct must implement the `BlockHandleModule` interface defined [here](../modules/block.go).
For clarity, here is the definition:

```go
// Module represents a module used to index a blockchain.
type Module interface {
	// GetName returns the name that identifies the module.
	GetName() string
}

// BlockHandleModule represents a module that indexes data by extracting it from a block.
type BlockHandleModule interface {
	Module
	// HandleBlock processes the provided block.
	HandleBlock(ctx context.Context, block types.Block) error
}
```

While implementing your parsing logic, you’ll notice that `HandleBlock` receives a generic `Block`, which lacks chain-specific fields.
To work with your custom block type, you must cast it explicitly:

```go
func (m *MyModule) HandleBlock(ctx context.Context, block types.Block) error {
	myBlock, ok := block.(*mytypes.Block)
	if !ok {
		// This block type is not handled by this module; ignore it.
		return nil
	}

	// Proceed with processing the block.
}
```

---

## Define a Custom `BlockHandleModule` Interface

If you prefer not to cast the block type manually every time, you can define a custom `BlockHandleModule` interface that directly receives your block type:

```go
type MyCustomBlockHandleModule interface {
	Module
	HandleBlock(ctx context.Context, block *mytypes.Block) error
}
```

Then, create an adapter that performs the type assertion and delegates the call to your custom module:

```go
type BlockHandleAdapter struct {
	handler MyCustomBlockHandleModule
}

var _ modules.BlockHandleModule = &BlockHandleAdapter{}

func NewBlockHandleAdapter(handler MyCustomBlockHandleModule) *BlockHandleAdapter {
	return &BlockHandleAdapter{
		handler: handler,
	}
}

func (b *BlockHandleAdapter) GetName() string {
	return b.handler.GetName()
}

func (b *BlockHandleAdapter) HandleBlock(ctx context.Context, block types.Block) error {
	customBlock, ok := block.(*mytypes.Block)
	if !ok {
		// Not the expected block type, skip processing.
		return nil
	}

	return b.handler.HandleBlock(ctx, customBlock)
}
```

This approach improves type safety and reduces repetitive casting logic in your indexing modules.

---

## Register Your Module

To make your module available for use by `Indexer`s, you must register it with the `ModulesManager` by calling the `RegisterModule(moduleName string, builder Builder)` method.

This function requires:

* `moduleName`: a unique string identifier for your module.
* `builder`: a function that builds and returns an instance of your module.

The builder function must follow this signature:

```go
// Builder is a function that constructs a Module instance.
// It receives:
// - ctx: the context for initialization
// - database: the database instance used by the indexer
// - node: the blockchain node used by the indexer
// - rawConfig: raw YAML configuration specific to the module
type Builder func(ctx context.Context, database database.Database, node node.Node, rawConfig []byte) (modules.Module, error)
```

### Purpose of `moduleName`

The `moduleName` serves two key roles:

1. **Configuration Binding**
   It links your module to indexer configurations. 
   When an `Indexer` specifies a module by name (see [Indexer Configuration](./config_structure.md#indexers)), 
   the `ModulesManager` uses the `moduleName` to locate the correct builder function.

2. **Module Configuration Association**
   In the global [modules configuration section](./config_structure.md#modules), 
   `moduleName` is used as the key under which the module's configuration is defined.
   These YAML values are passed to your builder and can be unmarshaled into your custom config struct.

By following this registration pattern, your custom module becomes fully integrated and configurable via the library's YAML-based configuration system.

Here’s an improved and polished version of your **"Register your adapter"** section with clearer grammar, structure, and explanation:

### Register Your Adapter

If you've implemented a custom module interface along with an adapter 
(as discussed in the previous section), you should return the adapter from your module builder function.

In the builder, construct your custom module and wrap it inside the adapter before returning.
Here's an example using a module for Cosmos-SDK based blockchains:

```go
func ExampleBlockBuilder(ctx context.Context, database database.Database, node node.Node, cfg []byte) (modules.Module, error) {
	indexerCtx := indexertypes.GetIndexerContext(ctx)

	return cosmosmodules.NewBlockHandleAdapter(&ExampleModule{
		logger: indexerCtx.Logger.With().Str("module", "example").Logger(),
	}), nil
}
```

In this example:

* `ExampleModule` is your custom module implementing the chain-specific handler interface.
* `NewBlockHandleAdapter` wraps the custom module to provide a generic `BlockHandleModule` interface implementation.
* The adapter ensures that the `Indexer` can invoke the handler using the standard interface while allowing you to work directly with your custom block types internally.

This pattern keeps your indexing logic clean and decoupled from the generic interface layer, while still integrating seamlessly with the overall system.

