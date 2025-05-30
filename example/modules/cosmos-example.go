package modules

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/milkyway-labs/flux/cosmos/types"
	"github.com/milkyway-labs/flux/database"
	"github.com/milkyway-labs/flux/modules"
	"github.com/milkyway-labs/flux/modules/adapter"
	"github.com/milkyway-labs/flux/node"
	indexertypes "github.com/milkyway-labs/flux/types"
)

const CosmosExampleModuleName = "cosmos-example"

func CosmosExampleBlockBuilder(ctx context.Context, database database.Database, node node.Node, cfg []byte) (modules.Module, error) {
	indexerCtx := indexertypes.GetIndexerContext(ctx)
	return adapter.NewBlockHandleAdapter(&CosmosExampleModule{
		logger: indexerCtx.Logger.With().Str("module", CosmosExampleModuleName).Logger(),
	}), nil
}

var _ adapter.BlockHandleModule[*types.Block] = &CosmosExampleModule{}

type CosmosExampleModule struct {
	logger zerolog.Logger
}

// GetName implements modules.BlockHandleModule.
func (e *CosmosExampleModule) GetName() string {
	return CosmosExampleModuleName
}

// HandleBlock implements modules.BlockHandleModule.
func (e *CosmosExampleModule) HandleBlock(_ context.Context, block *types.Block) error {
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
					Msg("got transfer event")
			}
		}
	}

	e.logger.Info().Uint64("height", uint64(block.GetHeight())).Msg("handled block")

	return nil
}
