package modules

import "context"

// Module represent a module used to index a block chain.
type Module interface {
	// GetName gets the name that identifies the module
	GetName() string
}

// IndexerStartHook represents a hook that is called when the indexer with the
// module is started.
type IndexerStartHook interface {
	OnIndexerStart(ctx context.Context) error
}
