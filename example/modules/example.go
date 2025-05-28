package modules

import (
	"context"

	"github.com/rs/zerolog"

	cosmosmodules "github.com/milkyway-labs/chain-indexer/cosmos/modules"
	"github.com/milkyway-labs/chain-indexer/cosmos/types"
	"github.com/milkyway-labs/chain-indexer/database"
	"github.com/milkyway-labs/chain-indexer/modules"
	"github.com/milkyway-labs/chain-indexer/node"
	indexertypes "github.com/milkyway-labs/chain-indexer/types"
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
	// for _, tx := range block.Txs {
	// 	events := tx.Events.GetEventsWithType("transfer")
	// 	for _, transferEvent := range events {
	// 		from, hasFrom := transferEvent.FindAttribute("sender")
	// 		to, hasTo := transferEvent.FindAttribute("recipient")
	// 		amount, hasAmount := transferEvent.FindAttribute("amount")
	// 		if hasFrom && hasTo && hasAmount {
	// 			// e.logger.Info().
	// 			// 	Str("from", from.Value).
	// 			// 	Str("to", to.Value).
	// 			// 	Str("amount", amount.Value).
	// 			// 	Msg("go transfer event")
	// 		}
	// 	}
	// }
	//
	// e.logger.Info().Uint64("height", uint64(block.GetHeight())).Msg("handled block")

	return nil
}
