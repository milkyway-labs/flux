package modules

import (
	"context"

	"github.com/milkyway-labs/chain-indexer/types"
)

// TxHandleModule represent a module that index data by extracting them from
// a tx.
type TxHandleModule interface {
	Module
	// HandleTx process the provided tx that is included in the provided block.
	HandleTx(ctx context.Context, block types.Block, tx types.Tx) error
}
