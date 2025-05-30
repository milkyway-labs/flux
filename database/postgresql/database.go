package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"github.com/milkyway-labs/flux/database"
	"github.com/milkyway-labs/flux/types"
)

// type check to ensure interface is properly implemented
var _ database.Database = &Database{}

// Database defines a wrapper around a SQL database and implements functionality
// for data aggregation and exporting.
type Database struct {
	Logger zerolog.Logger
	Cfg    *Config
	SQL    *sqlx.DB
}

type BlockRow struct {
	ChainID   string       `db:"chain_id"`
	Height    types.Height `db:"height"`
	Timestamp time.Time    `db:"timestamp"`
}

func NewDatabase(logger zerolog.Logger, cfg *Config) (*Database, error) {
	postgresDB, err := sqlx.Open("postgres", cfg.URL)
	if err != nil {
		return nil, err
	}

	return &Database{
		Logger: logger.With().Str("component", "database").Logger(),
		Cfg:    cfg,
		SQL:    postgresDB,
	}, nil
}

// GetLowestBlock implements database.Database.
func (db *Database) GetLowestBlock(chainID string) (*types.Height, error) {
	stmt := `SELECT height FROM blocks WHERE chain_id = $1 ORDER BY height ASC LIMIT 1`

	var height types.Height
	err := db.SQL.QueryRow(stmt, chainID).Scan(&height)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &height, nil
}

// GetMissingBlocks implements database.Database.
func (db *Database) GetMissingBlocks(chainID string, from types.Height, to types.Height) ([]types.Height, error) {
	if from > to {
		return nil, fmt.Errorf("invalid range, from(%d) must not be greater than to(%d)", from, to)
	}

	var result []types.Height
	stmt := `SELECT generate_series($1::int,$2::int) EXCEPT SELECT height FROM blocks WHERE chain_id = $3 ORDER BY 1`

	err := db.SQL.Select(&result, stmt, from, to, chainID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SaveIndexedBlock implements database.Database.
func (db *Database) SaveIndexedBlock(chainID string, height types.Height, timestamp time.Time) error {
	stmt := `
INSERT INTO blocks (chain_id, height, timestamp)
VALUES ($1, $2, $3)
ON CONFLICT ON CONSTRAINT unique_chain_block DO UPDATE
	SET timestamp = excluded.timestamp
		`

	_, err := db.SQL.Exec(stmt,
		chainID,
		height,
		timestamp.UTC(),
	)
	return err
}
