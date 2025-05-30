package manager

import (
	"context"
	"fmt"

	"github.com/milkyway-labs/flux/database"
)

// DatabasesManager handle the construction of the Database instances that can
// be used by an indexer to store the indexed data.
type DatabasesManager struct {
	registered map[string]Builder
}

func NewDatabasesManager() *DatabasesManager {
	return &DatabasesManager{
		registered: make(map[string]Builder),
	}
}

// RegisterDatabase register a new database type that can be used by an indexer to
// store the indexed data.
func (mm *DatabasesManager) RegisterDatabase(dbType string, builder Builder) *DatabasesManager {
	mm.registered[dbType] = builder
	return mm
}

// GetDatabase builds an return a Database instance having the requested type.
func (mm *DatabasesManager) GetDatabase(
	ctx context.Context,
	dbType string,
	databaseID string,
	cfg []byte,
) (database.Database, error) {
	// Get the database builder
	builder, found := mm.registered[dbType]
	if !found {
		return nil, fmt.Errorf("can't find builder for db `%s`", dbType)
	}

	return builder(ctx, databaseID, cfg)
}
