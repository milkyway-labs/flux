# Cosmos-SDK Chains

Here is the code that provides indexing support for Cosmos-SDK-based blockchains.

## Registration

To enable indexing of Cosmos-SDK-based blockchains, you need to register the 
Cosmos `Node` implementation that can fetch `Block`s from a node. 
You can do that with the following code:

```go
import (
	cosmosrpc "github.com/milkyway-labs/flux/cosmos/node/rpc"
)

// Register the Cosmos Node in the NodesManager used by the IndexerBuilder
nodesManager.RegisterNode(cosmosrpc.NodeType, cosmosrpc.NodeBuilder)
```

### Configuration

Below is an example of a valid Cosmos node configuration:

```yaml
type: "cosmos-rpc"
url: "https://rpc.chain.zone"
request_timeout: "10s"
```

**Fields:**

* `type`: Specifies the node type so the library can instantiate the correct `Node` implementation.
* `url`: The node's RPC URL.
* `request_timeout`: The amount of time the client will wait for a response from the node 
before considering the request failed. Defaults to `10s`.
* `tx_events_from_log_until_height`: Specifies the height until which the `tx.log` field will be 
used to extract transaction events. After this height, the `tx.events` field will 
be used instead. If this field is undefined, `tx.events` will always be used.
* `decode_block_event_attributes_until_height`: Specifies the height until which block events are 
treated as base64-encoded and need to be decoded. If this field is undefined, block events will not be decoded.

## Cosmos Modules

To create a `Module` capable of indexing Cosmos-SDK-based blockchains, 
define a `struct` that implements either the `BlockHandleModule` or `TxHandleModule` from the `github.com/milkyway-labs/flux/cosmos/modules` package.

Below is an example of a `BlockHandleModule` that logs transfer actions:

```go
package modules

import (
	"context"

	"github.com/rs/zerolog"
	cosmosmodules "github.com/milkyway-labs/flux/cosmos/modules"
	"github.com/milkyway-labs/flux/cosmos/types"
	"github.com/milkyway-labs/flux/database"
	"github.com/milkyway-labs/flux/modules"
	"github.com/milkyway-labs/flux/node"
	indexertypes "github.com/milkyway-labs/flux/types"
)

var _ cosmosmodules.BlockHandleModule = &ExampleModule{}

type ExampleModule struct {
	logger zerolog.Logger
}

func ExampleBlockBuilder(ctx context.Context, database database.Database, node node.Node, cfg []byte) (modules.Module, error) {
	indexerCtx := indexertypes.GetIndexerContext(ctx)
	return cosmosmodules.NewBlockHandleAdapter(&ExampleModule{
		logger: indexerCtx.Logger.With().Str("module", "example").Logger(),
	}), nil
}

// GetName implements modules.BlockHandleModule.
func (e *ExampleModule) GetName() string {
	return "example"
}

// HandleBlock implements modules.BlockHandleModule.
func (e *ExampleModule) HandleBlock(ctx context.Context, block *types.Block) error {
	for _, tx := range block.Txs {
		for _, transferEvent := range tx.Events.FindEventsWithType("transfer") {
			from, hasFrom := transferEvent.FindAttribute("sender")
			to, hasTo := transferEvent.FindAttribute("recipient")
			amount, hasAmount := transferEvent.FindAttribute("amount")
			if hasFrom && hasTo && hasAmount {
				e.logger.Info().
					Str("from", from.Value).
					Str("to", to.Value).
					Str("amount", amount.Value).
					Msg("go transfer event")
			}
		}
	}

	e.logger.Info().Uint64("height", uint64(block.GetHeight())).Msg("handled block")

	return nil
}
```

### Registration

After creating your custom `Module`, you must register it to be used by an `Indexer`.
For Cosmos-SDK chains, we provide `BlockHandleAdapter` and `TxHandleAdapter` components 
that enable the registration of Cosmos-specific modules and ensure they are called only when indexing a `Block` produced by a Cosmos-SDK `Node`.  

Below is an example that shows how to use the `BlockHandleAdapter` to register a Cosmos module:

```go
import (
	"context"

	cosmosmodules "github.com/milkyway-labs/flux/cosmos/modules"
	"github.com/milkyway-labs/flux/database"
	"github.com/milkyway-labs/flux/node"
	indexertypes "github.com/milkyway-labs/flux/types"
)

func ExampleBlockBuilder(ctx context.Context, database database.Database, node node.Node, cfg []byte) (modules.Module, error) {
	indexerCtx := indexertypes.GetIndexerContext(ctx)
	return cosmosmodules.NewBlockHandleAdapter(&ExampleModule{
		logger: indexerCtx.Logger.With().Str("module", "example").Logger(),
	}), nil
}


// Register the example module with the ModulesManager used by the IndexerBuilder
modulesManager.RegisterModule("example", ExampleBlockBuilder)
```

