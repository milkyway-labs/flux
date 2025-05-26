package database

import (
	"time"

	"github.com/milkyway-labs/chain-indexer/types"
)

// Database represents a database used by the indexer to store the indexing state.
type Database interface {
	// GetLowestBlock retrieves the height of the lowest indexed block for the provided
	// chainID.
	// If no blocks have been indexed for the specified chain, a nil height is returned.
	GetLowestBlock(chainID string) (*types.Height, error)
	// GetMissingBlocks retrieves the blocks that need to be indexed from the chain
	// with the provided chainID in the provided bock range.
	// A block is considered missing if it has not been indexed yet
	// or if a previous indexing operation failed.
	GetMissingBlocks(chainID string, from types.Height, to types.Height) ([]types.Height, error)
	// Stores in the database that the given height for the chain with the provided ID
	// has been indexed.
	SaveIndexedBlock(chainID string, height types.Height, timestamp time.Time) error
}
