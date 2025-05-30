package modules

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/milkyway-labs/flux/database"
	evmtypes "github.com/milkyway-labs/flux/evm/types"
	"github.com/milkyway-labs/flux/modules"
	"github.com/milkyway-labs/flux/modules/adapter"
	"github.com/milkyway-labs/flux/node"
	indexertypes "github.com/milkyway-labs/flux/types"
)

const (
	EVMExampleModuleName = "evm-example"
	ERC20TransferTopic   = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

func EVMExampleBlockBuilder(ctx context.Context, database database.Database, node node.Node, cfg []byte) (modules.Module, error) {
	indexerCtx := indexertypes.GetIndexerContext(ctx)
	return adapter.NewBlockHandleAdapter(&EVMExampleModule{
		logger: indexerCtx.Logger.With().Str("module", EVMExampleModuleName).Logger(),
	}), nil
}

var _ adapter.BlockHandleModule[*evmtypes.Block] = &EVMExampleModule{}

type EVMExampleModule struct {
	logger zerolog.Logger
}

// GetName implements adpter.BlockHandleModule.
func (e *EVMExampleModule) GetName() string {
	return EVMExampleModuleName
}

// HandleBlock implements adpter.BlockHandleModule.
func (e *EVMExampleModule) HandleBlock(ctx context.Context, block *evmtypes.Block) error {
	for _, log := range block.Logs {
		topic := log.Topics[0].Hex()
		if topic == ERC20TransferTopic && len(log.Topics) == 3 {
			e.logger.Info().
				Str("address", log.Address.NormalizedHex()).
				Str("from", log.Topics[1].NormalizedHex()).
				Str("to", log.Topics[2].NormalizedHex()).
				Str("amount", log.Data.Int().String()).
				Msg("ERC20 transfer event")
		}
	}

	e.logger.Info().Uint64("height", uint64(block.GetHeight())).Msg("handled block")

	return nil
}
