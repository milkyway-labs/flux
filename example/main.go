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
	postgresql.AddPostgressDatabaseSupport(ctx.DatabasesManager)
	rpc.AddCosmosRPCNodeSupport(ctx.NodesManager)
	ctx.ModulesManager.RegisterModule("example", modules.ExampleBlockBuilder)

	err := cli.NewSimpleCLI(ctx).Execute()
	if err != nil {
		panic(err)
	}
}
