# EVM Chains

Here is the code that provides indexing support for EVM blockchains.

## Registration

To enable indexing of EVM blockchains, you need to register the 
EVM `Node` implementation that can fetch `Block`s from a node. 
You can do that with the following code:

```go
import (
	evmrpc "github.com/milkyway-labs/chain-indexer/evm/node/rpc"
)

// Register the Cosmos Node in the NodesManager used by the IndexerBuilder
nodesManager.RegisterNode(evmrpc.NodeType, evmrpc.NodeBuilder)
```

### Configuration

Below is an example of a valid EVM node configuration:

```yaml
type: "evm-rpc"
url: "https://rpc.chain.zone"
request_timeout: "10s"
```

**Fields:**

* `type`: Specifies the node type so the library can instantiate the correct `Node` implementation.
* `url`: The node's RPC URL.
* `request_timeout`: The amount of time the client will wait for a response from the node 
before considering the request failed. Defaults to `10s`.

## EVM Modules

To create a `Module` capable of indexing Cosmos-SDK-based blockchains, 
define a `struct` that implements either the `BlockHandleModule[*evmtypes.Block]` or 
`BlockHandleModule[*evmtypes.Block, *evmtypes.Tx]` from the `github.com/milkyway-labs/chain-indexer/modules/adapter` package.

Below is an example of a `BlockHandleModule` that logs transfer actions:

```go
package modules

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/milkyway-labs/chain-indexer/modules/adapter"
	evmtypes "github.com/milkyway-labs/chain-indexer/evm/types"
	"github.com/milkyway-labs/chain-indexer/database"
	"github.com/milkyway-labs/chain-indexer/modules"
	"github.com/milkyway-labs/chain-indexer/node"
	indexertypes "github.com/milkyway-labs/chain-indexer/types"
)

var _ adapter.BlockHandleModule[*evmtypes.Block] = &ExampleModule{}

const Name = "example"

type ExampleModule struct {
	logger zerolog.Logger
}

func Builder(ctx context.Context, database database.Database, node node.Node, cfg []byte) (modules.Module, error) {
	indexerCtx := indexertypes.GetIndexerContext(ctx)
	return adapter.NewBlockHandleAdapter(&ExampleModule{
		logger: indexerCtx.Logger.With().Str("module", "example").Logger(),
	}), nil
}

// GetName implements modules.BlockHandleModule.
func (e *ExampleModule) GetName() string {
	return Name
}

// HandleBlock implements modules.BlockHandleModule.
func (e *ExampleModule) HandleBlock(ctx context.Context, block *evmtypes.Block) error {
	e.logger.Info().Uint64("height", uint64(block.GetHeight())).Msg("handled block")

    // Your logic here

	return nil
}
```

### Registration

After creating your custom `Module`, you must register it to be used by an `Indexer`.

Below is an example that shows how to use the `BlockHandleAdapter` to register a Cosmos module:

```go
import (
	"context"

	examplemodule "github.com/your-nae/modules/example"
)

// Register the example module with the ModulesManager used by the IndexerBuilder
modulesManager.RegisterModule(examplemodule.Name, examplemodule.Builder)
```

