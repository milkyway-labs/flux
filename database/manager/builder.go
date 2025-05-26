package manager

import (
	"context"

	"github.com/milkyway-labs/chain-indexer/database"
)

// Builder represents a function that can be used to build a Database instance,
// this functions will receive the raw yaml bytes and an ID that identifies the database.
type Builder func(ctx context.Context, id string, rawConfig []byte) (database.Database, error)
