package main

import (
	"github.com/milkyway-labs/flux/cli"
	"github.com/milkyway-labs/flux/cli/types"
	cosmosrpc "github.com/milkyway-labs/flux/cosmos/node/rpc"
	"github.com/milkyway-labs/flux/database/postgresql"
	evmrpc "github.com/milkyway-labs/flux/evm/node/rpc"
	"github.com/milkyway-labs/flux/example/modules"
)

func main() {
	ctx := types.NewCliContext("example")
	// Database types
	ctx.DatabasesManager.RegisterDatabase(postgresql.DatabaseType, postgresql.DatabaseBuilder)

	// Nodes types
	ctx.NodesManager.RegisterNode(cosmosrpc.NodeType, cosmosrpc.NodeBuilder)
	ctx.NodesManager.RegisterNode(evmrpc.NodeType, evmrpc.NodeBuilder)

	// Modules
	ctx.ModulesManager.RegisterModule(modules.CosmosExampleModuleName, modules.CosmosExampleBlockBuilder)
	ctx.ModulesManager.RegisterModule(modules.EVMExampleModuleName, modules.EVMExampleBlockBuilder)

	err := cli.NewDefaultIndexerCLI(ctx).Execute()
	if err != nil {
		panic(err)
	}
}
