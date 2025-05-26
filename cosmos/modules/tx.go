package modules

import (
	"context"

	cosmostypes "github.com/milkyway-labs/chain-indexer/cosmos/types"
	"github.com/milkyway-labs/chain-indexer/modules"
)

// TxHandleModule represent a module that index data by extracting them from
// a tx.
type TxHandleModule interface {
	modules.Module
	// HandleTx process the provided tx that is included in the provided block.
	HandleTx(ctx context.Context, block *cosmostypes.Block, tx *cosmostypes.Tx) error
}
