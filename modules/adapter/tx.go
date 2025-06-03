package adapter

import (
	"context"

	"github.com/milkyway-labs/flux/modules"
	"github.com/milkyway-labs/flux/types"
)

type TxHandleModule[B types.Block, T types.Tx] interface {
	modules.Module
	// HandleTx process the provided tx that is included in the provided block.
	HandleTx(ctx context.Context, block B, tx T) error
}

// ----------------------------------------------------------------------------
// ---- Tx module adapter
// ----------------------------------------------------------------------------

type TxHandleAdapter[B types.Block, T types.Tx] struct {
	handler TxHandleModule[B, T]
}

func NewTxHandleAdapter[B types.Block, T types.Tx](handler TxHandleModule[B, T]) modules.TxHandleModule {
	return &TxHandleAdapter[B, T]{
		handler: handler,
	}
}

// GetName implements indexer.TxHandleModule.
func (b *TxHandleAdapter[B, T]) GetName() string {
	return b.handler.GetName()
}

// HandleTx implements indexer.TxHandleModule.
func (b *TxHandleAdapter[B, T]) HandleTx(ctx context.Context, block types.Block, tx types.Tx) error {
	castedBlock, ok := block.(B)
	if !ok {
		return nil
	}
	castedTx, ok := tx.(T)
	if !ok {
		return nil
	}

	return b.handler.HandleTx(ctx, castedBlock, castedTx)
}
