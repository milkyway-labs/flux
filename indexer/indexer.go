package indexer

import (
	"context"
	"fmt"
	"sync"

	log "github.com/rs/zerolog"

	"github.com/milkyway-labs/flux/database"
	"github.com/milkyway-labs/flux/modules"
	"github.com/milkyway-labs/flux/node"
	"github.com/milkyway-labs/flux/types"
)

type Indexer struct {
	// Indexer's configuration
	cfg *types.IndexerConfig

	// Logger used by the indexer
	log log.Logger

	// Database used by the indexer to store its state.
	db database.Database

	// The node from which the indexer fetches the blocks to index.
	node node.Node

	// Channel used by the workers to retrieve the height of the
	// blocks to index.
	heightsQueue *Queue[IndexerHeight]

	// List of modules that will be used by the indexer to index data from
	// the chain.
	modules []modules.Module

	// Instance of HeightProducer that will provide the blocks to parse.
	heightProducer HeightProducer
}

func NewIndexer(
	cfg *types.IndexerConfig,
	log log.Logger,
	db database.Database,
	node node.Node,
	modules []modules.Module,
) Indexer {
	logger := log.With().
		Str("indexer", cfg.Name).
		Str("chain-id", node.GetChainID()).
		Logger()
	return Indexer{
		cfg:          cfg,
		log:          logger,
		db:           db,
		node:         node,
		heightsQueue: NewQueue[IndexerHeight](cfg.HeightQueueSize),
		modules:      modules,
	}
}

// GetName get the name that identifies the indexer.
func (i *Indexer) GetName() string {
	return i.cfg.Name
}

// Start starts the indexer.
func (i *Indexer) Start(ctx context.Context, wg *sync.WaitGroup) error {
	heightProducer := i.heightProducer

	// If we don't have a height producer, we build the default one.
	if heightProducer == nil {
		producer, err := i.buildDefaultHeightProducer(ctx)
		if err != nil {
			return fmt.Errorf("build default height producer: %w", err)
		}
		heightProducer = producer
	}

	// Start the worker that produces the heights to be fetched by the workers.
	wg.Add(1)
	go i.enqueueHeightsLoop(ctx, wg, heightProducer)

	// Starts the indexing workers
	for index := int64(0); index < int64(i.cfg.Workers); index++ {
		worker := NewWorker(i.cfg, i.log, i.heightsQueue, i.db, i.node, i.modules)
		worker.Start(ctx, wg)
	}

	return nil
}

// WithCustomHeightProducer allows to define a custom HeightProducer that provides
// the heights to parse.
func (i *Indexer) WithCustomHeightProducer(producer HeightProducer) *Indexer {
	i.heightProducer = producer
	return i
}

// buildDefaultHeightProducer builds the default height producer, this producer
// will produce the height of the un-indexed blocks from first indexed block or the configured start height
// in case is the first start to the current node height and starts monitor the node for new blocks.
func (i *Indexer) buildDefaultHeightProducer(ctx context.Context) (HeightProducer, error) {
	currentNodeHeight, err := i.node.GetCurrentHeight(ctx)
	if err != nil {
		return nil, fmt.Errorf("get current node height: %w", err)
	}

	var missingBlockStartHeight types.Height
	if i.cfg.StartHeight != nil {
		// Start from the height specified in the config
		missingBlockStartHeight = *i.cfg.StartHeight
	} else {
		// Get the lowest indexed block.
		lowestAvailableBlock, err := i.db.GetLowestBlock(i.node.GetChainID())
		if err != nil {
			return nil, fmt.Errorf("get lowest block %w", err)
		}

		// In case we have an indexed block and is lower then the current node
		// height check the missing block from this height
		if lowestAvailableBlock != nil && *lowestAvailableBlock < currentNodeHeight {
			missingBlockStartHeight = *lowestAvailableBlock
		} else {
			// We don't have any indexed block or the current node height is lower
			// than the lowest indexed block, which is wired. In those cases
			// start looking for un-indexed block from the current node height.
			missingBlockStartHeight = currentNodeHeight
		}
	}

	var missingBlocks []types.Height
	if i.cfg.ForceReparseOldBlocks && i.cfg.StartHeight != nil {
		for height := *i.cfg.StartHeight; height <= currentNodeHeight; height++ {
			missingBlocks = append(missingBlocks, height)
		}
	} else {
		// Get the blocks that are missing and we need to index
		missingBlocksFromDB, err := i.db.GetMissingBlocks(i.node.GetChainID(), missingBlockStartHeight, currentNodeHeight-1)
		if err != nil {
			return nil, fmt.Errorf("get missing blocks: %w", err)
		}
		missingBlocks = missingBlocksFromDB
	}

	return NewCombinedHeightProducer(
		NewListHeightProducer(missingBlocks),
		NewNodeHeightProducer(i.log, i.node, i.cfg.NodePollingInterval, currentNodeHeight),
	), nil
}

func (i *Indexer) enqueueHeightsLoop(
	ctx context.Context,
	wg *sync.WaitGroup,
	heightProducer HeightProducer,
) {
	defer func() {
		i.heightsQueue.Close()
		wg.Done()
	}()

	heightProducer.EnqueueHeights(ctx, i.heightsQueue)
}
