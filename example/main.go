package main

import (
	"github.com/milkyway-labs/flux/cli"
	"github.com/milkyway-labs/flux/cli/types"
	"github.com/milkyway-labs/flux/cosmos/node/rpc"
	"github.com/milkyway-labs/flux/database/postgresql"
	"github.com/milkyway-labs/flux/example/modules"
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
