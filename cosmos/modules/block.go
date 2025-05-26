package modules

import (
	"context"

	cosmostypes "github.com/milkyway-labs/chain-indexer/cosmos/types"
	"github.com/milkyway-labs/chain-indexer/modules"
)

// BlockHandleModule represent a module that index data by extracting them from
// a block.
type BlockHandleModule interface {
	modules.Module
	// HandleBlock process the provided block.
	HandleBlock(ctx context.Context, block *cosmostypes.Block) error
}
