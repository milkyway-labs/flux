package manager

import (
	"context"

	"github.com/milkyway-labs/chain-indexer/database"
	"github.com/milkyway-labs/chain-indexer/modules"
	"github.com/milkyway-labs/chain-indexer/node"
)

// Builder represents a function that can be used to build a Module instance,
// this functions will receive the module's config as yaml bytes, the database and node used by the indexer
// where the module will be used.
type Builder func(ctx context.Context, datdatabase database.Database, node node.Node, rawConfig []byte) (modules.Module, error)
