package node

import (
	"context"

	"github.com/milkyway-labs/chain-indexer/types"
)

// Node represents a generic block chain node that can be queried to
// obtain Blocks
type Node interface {
	// GetChainID gets the ID that identifies the block chain that is being queried.
	GetChainID() string
	// GetBlock queries the node to get the block produced at the provided
	// height.
	GetBlock(context context.Context, height types.Height) (types.Block, error)
	// GetLowestHeight gets the lowest height that can be queried from the node.
	GetLowestHeight(context context.Context) (types.Height, error)
	// GetCurrentHeight gets the current node height.
	GetCurrentHeight(context context.Context) (types.Height, error)
}
