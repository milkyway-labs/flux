package main

import (
	"github.com/milkyway-labs/chain-indexer/cli"
	"github.com/milkyway-labs/chain-indexer/cli/types"
	"github.com/milkyway-labs/chain-indexer/cosmos/node/rpc"
	"github.com/milkyway-labs/chain-indexer/database/postgresql"
	"github.com/milkyway-labs/chain-indexer/example/modules"
)

func main() {
	ctx := types.NewCliContext("example")
	// Database types
	ctx.DatabasesManager.RegisterDatabase(postgresql.DatabaseType, postgresql.DatabaseBuilder)

	// Nodes types
	ctx.NodesManager.RegisterNode(rpc.NodeType, rpc.NodeBuilder)

	// Modules
	ctx.ModulesManager.RegisterModule("example", modules.ExampleBlockBuilder)

	err := cli.NewDefaultIndexerCLI(ctx).Execute()
	if err != nil {
		panic(err)
	}
}
