package manager

import (
	"context"

	"github.com/milkyway-labs/chain-indexer/node"
)

// Builder represents a function that can be used to build a Node instance,
// this functions will receive the raw yaml bytes and an ID that identifies the node.
type Builder func(ctx context.Context, id string, rawConfig []byte) (node.Node, error)
