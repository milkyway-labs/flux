# Support new database

This document explains how to support a new database backend that can be used 
by the `Indexer` or a `Module` to store the indexed data. 

---


## Implement the Database interface

To create a database instance that can be used by the `Indexer`,
you need to define a structure that implements the `Database` interface.
Here is its definition:

```go
// Database represents a database used by the indexer to store indexing state.
type Database interface {
	// GetLowestBlock retrieves the height of the lowest indexed block for the provided
	// chainID.
	// If no blocks have been indexed for the specified chain, a nil height is returned.
	GetLowestBlock(chainID string) (*types.Height, error)

	// GetMissingBlocks retrieves the blocks that need to be indexed from the chain
	// with the provided chainID, within the specified block range.
	// A block is considered missing if it has not been indexed yet
	// or if a previous indexing operation failed.
	GetMissingBlocks(chainID string, from types.Height, to types.Height) ([]types.Height, error)

	// SaveIndexedBlock records in the database that the given height for the specified chain
	// has been successfully indexed.
	SaveIndexedBlock(chainID string, height types.Height, timestamp time.Time) error
}
```

Once you have implemented this interface, your `Database` instance can be used by the
indexer to store indexing state. It can also be extended to store module-specific data
when building an indexer for a particular use case.

For a reference implementation you can look at the PostgresSQL implementation [here](../database/postgresql/database.go)

## Register your Database type

To register your custom `Database` and allow the library to build an instance of it,
you must provide a `Builder` function with the following signature:

```go
// Builder represents a function that can be used to build a Database instance,
// this functions will receive the raw yaml bytes and an ID that identifies the database.
type Builder func(ctx context.Context, id string, rawConfig []byte) (database.Database, error)
```

This allows the `DatabasesManager` to construct your database instance when an `Indexer` requires it.

To register your node, call the `DatabasesManager`'s `RegisterDatabase(type string, builder Builder)` function,
providing the database type identifier and the corresponding `Builder` function.
The `type` is used in the configuration file to determine which builder should be used to initialize the given `Database` type.

For example, if you’ve created an `MariaDBDatabase` for MariaDB registered it with `type` set to `"mariadb"`,
the library will try to build an `MariaDBDatabase` whenever it finds an entry like the following in the [`databases` section](./config_structure.md#databases) of the configuration:

```yaml
type: "mariadb"
# other values that will be passed as configuration during build time
```

For a reference implementation, see the `Builder` function for PostgresSQL `Database` [here](../database/postgresql/database.go).


