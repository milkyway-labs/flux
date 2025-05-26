package postgresql

import (
	"context"
	"fmt"

	"github.com/milkyway-labs/chain-indexer/database"
	"github.com/milkyway-labs/chain-indexer/database/manager"
	"github.com/milkyway-labs/chain-indexer/types"
	"gopkg.in/yaml.v3"
)

const PostgresDatabaseType = "postgres"

func PostgresDatabaseBuilder(
	ctx context.Context,
	id string,
	rawConfig []byte,
) (database.Database, error) {
	var config Config
	err := yaml.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal postgres db config %w", err)
	}

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid postgres db config %w", err)
	}

	indexerCtx := types.GetIndexerContext(ctx)

	return NewDatabase(indexerCtx.Logger, &config)
}

func AddPostgressDatabaseSupport(manager *manager.DatabasesManager) *manager.DatabasesManager {
	return manager.RegisterDatabase(PostgresDatabaseType, PostgresDatabaseBuilder)
}
