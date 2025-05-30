package postgresql_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/flux/database/postgresql"
)

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

type DbTestSuite struct {
	suite.Suite

	database *postgresql.Database
}

func (suite *DbTestSuite) SetupTest() {
	// Build the database config
	dbCfg := postgresql.DefaultConfig().
		WithURL("postgres://milkyway:password@localhost:6432/milkyway?sslmode=disable&search_path=public")

	// Build the database
	parserDb, err := postgresql.NewDatabase(log.Logger, &dbCfg)
	suite.Require().NoError(err)

	// Delete the public schema
	_, err = parserDb.SQL.Exec(`DROP SCHEMA public CASCADE;`)
	suite.Require().NoError(err)

	// Create the schema
	_, err = parserDb.SQL.Exec(`CREATE SCHEMA public;`)
	suite.Require().NoError(err)

	dirPath := "schema"
	dir, err := os.ReadDir(dirPath)
	suite.Require().NoError(err)

	for _, fileInfo := range dir {
		if !strings.HasSuffix(fileInfo.Name(), ".sql") {
			continue
		}

		file, err := os.ReadFile(filepath.Join(dirPath, fileInfo.Name()))
		suite.Require().NoError(err)

		_, err = parserDb.SQL.Exec(string(file))
		suite.Require().NoError(err)
	}

	// Create the truncate function
	stmt := fmt.Sprintf(`
CREATE OR REPLACE FUNCTION truncate_tables(username IN VARCHAR) RETURNS void AS $$
DECLARE
    table_statements CURSOR FOR
        SELECT tablename FROM pg_tables
        WHERE tableowner = username AND schemaname = '%[1]s';
    sequence_statements CURSOR FOR
        SELECT sequence_name FROM information_schema.sequences
        WHERE sequence_schema = '%[1]s';
BEGIN
    FOR stmt IN table_statements LOOP
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.tablename) || ' CASCADE;';
    END LOOP;
    
    FOR seq IN sequence_statements LOOP
        EXECUTE 'ALTER SEQUENCE ' || quote_ident(seq.sequence_name) || ' RESTART WITH 1;';
    END LOOP;
END;
$$ LANGUAGE plpgsql;`, "public")
	_, err = parserDb.SQL.Exec(stmt)
	suite.Require().NoError(err)

	suite.database = parserDb
}
