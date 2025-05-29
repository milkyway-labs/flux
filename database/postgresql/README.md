# PostgreSQL Driver

This section provides the code for the PostgreSQL `Database` implementation.
It can be used by indexers to store indexed blocks, and by `Module`s to store the data they extract.

## Registration

To register this database type, use the following code:

```go
import (
	"github.com/milkyway-labs/chain-indexer/database/postgresql"
)

// Register the PostgreSQL driver with the DatabaseManager used by the IndexerBuilder
databaseManager.RegisterDatabase(postgresql.DatabaseType, postgresql.DatabaseBuilder)
```

### Configuration

Below is an example of a valid PostgreSQL database configuration:

```yaml
type: "postgres"
url: "postgres://milkyway:password@localhost:6432/milkyway?sslmode=disable&search_path=public"
partition_size: 100000
```

**Fields:**

* `type`: Specifies the database type so the library can instantiate the correct driver.
* `url`: The URI used to connect to the database.
* `partition_size`: The partition size used by the driver (default: 100,000).

