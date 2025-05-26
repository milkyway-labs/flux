package modules

import (
	"context"

	"github.com/milkyway-labs/chain-indexer/types"
)

// BlockHandleModule represent a module that index data by extracting them from
// a block.
type BlockHandleModule interface {
	Module
	// HandleBlock process the provided block.
	HandleBlock(ctx context.Context, block types.Block) error
}
