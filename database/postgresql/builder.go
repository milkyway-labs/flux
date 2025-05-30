package postgresql

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/milkyway-labs/flux/database"
	"github.com/milkyway-labs/flux/types"
)

const DatabaseType = "postgres"

func DatabaseBuilder(
	ctx context.Context,
	_ string,
	rawConfig []byte,
) (database.Database, error) {
	var config Config
	err := yaml.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal postgres db config: %w", err)
	}

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid postgres db config: %w", err)
	}

	indexerCtx := types.GetIndexerContext(ctx)

	return NewDatabase(indexerCtx.Logger, &config)
}
