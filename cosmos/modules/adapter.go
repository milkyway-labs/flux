package modules

import (
	"context"

	cosmostypes "github.com/milkyway-labs/chain-indexer/cosmos/types"
	"github.com/milkyway-labs/chain-indexer/modules"
	"github.com/milkyway-labs/chain-indexer/types"
)

// ----------------------------------------------------------------------------
// ---- Tx module adapter
// ----------------------------------------------------------------------------

type TxHandleAdapter struct {
	handler TxHandleModule
}

var _ modules.TxHandleModule = &TxHandleAdapter{}

func NewTxHandleAdapter(handler TxHandleModule) *TxHandleAdapter {
	return &TxHandleAdapter{
		handler: handler,
	}
}

// GetName implements indexer.TxHandleModule.
func (b *TxHandleAdapter) GetName() string {
	return b.handler.GetName()
}

// TxHandleModule implements indexer.TxHandleModule.
func (b *TxHandleAdapter) HandleTx(ctx context.Context, block types.Block, tx types.Tx) error {
	cosmosblock, ok := block.(*cosmostypes.Block)
	if !ok {
		return nil
	}
	cosmosTx, ok := tx.(*cosmostypes.Tx)
	if !ok {
		return nil
	}

	return b.handler.HandleTx(ctx, cosmosblock, cosmosTx)
}

// ----------------------------------------------------------------------------
// ---- Block module adapter
// ----------------------------------------------------------------------------

type BlockHandleAdapter struct {
	handler BlockHandleModule
}

var _ modules.BlockHandleModule = &BlockHandleAdapter{}

func NewBlockHandleAdapter(handler BlockHandleModule) *BlockHandleAdapter {
	return &BlockHandleAdapter{
		handler: handler,
	}
}

// GetName implements indexer.BlockHandleModule.
func (b *BlockHandleAdapter) GetName() string {
	return b.handler.GetName()
}

// HandleBlock implements indexer.BlockHandleModule.
func (b *BlockHandleAdapter) HandleBlock(ctx context.Context, block types.Block) error {
	cosmosblock, ok := block.(*cosmostypes.Block)
	if !ok {
		return nil
	}

	return b.handler.HandleBlock(ctx, cosmosblock)
}
