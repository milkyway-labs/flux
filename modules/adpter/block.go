package adpter

import (
	"context"

	"github.com/milkyway-labs/chain-indexer/modules"
	"github.com/milkyway-labs/chain-indexer/types"
)

// BlockHandleModule represent a module that index data by extracting them from
// a block.
type BlockHandleModule[B types.Block] interface {
	modules.Module
	// HandleBlock process the provided block.
	HandleBlock(ctx context.Context, block B) error
}

// ----------------------------------------------------------------------------
// ---- Block module adapter
// ----------------------------------------------------------------------------

type BlockHandleAdapter[B types.Block] struct {
	handler BlockHandleModule[B]
}

func NewBlockHandleAdapter[B types.Block](handler BlockHandleModule[B]) modules.BlockHandleModule {
	return &BlockHandleAdapter[B]{
		handler: handler,
	}
}

// GetName implements indexer.BlockHandleModule.
func (b *BlockHandleAdapter[B]) GetName() string {
	return b.handler.GetName()
}

// HandleBlock implements indexer.BlockHandleModule.
func (b *BlockHandleAdapter[B]) HandleBlock(ctx context.Context, block types.Block) error {
	cosmosblock, ok := block.(B)
	if !ok {
		return nil
	}

	return b.handler.HandleBlock(ctx, cosmosblock)
}
