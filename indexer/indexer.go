package indexer

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/rs/zerolog"

	"github.com/milkyway-labs/chain-indexer/database"
	"github.com/milkyway-labs/chain-indexer/modules"
	"github.com/milkyway-labs/chain-indexer/node"
	"github.com/milkyway-labs/chain-indexer/types"
	"github.com/milkyway-labs/chain-indexer/utils"
)

// indexHeight represents the struct used by the indexer to
// instruct a worker to fetch a block from the chain and keep track of
// the fetch attempts.
type indexHeight struct {
	Height   types.Height
	Attempts uint32
}

func newIndexHeight(height types.Height) indexHeight {
	return indexHeight{
		Height:   height,
		Attempts: 0,
	}
}

type Indexer struct {
	cfg *types.IndexerConfig
	log log.Logger

	// Database used by the indexer to store its state.
	db database.Database

	// The node from which the indexer fetches the blocks to index.
	node node.Node

	// Channel used by the workers to retrieve the height of the
	// blocks to index.
	HeightsQueue *Queue[indexHeight]

	// List of modules that will be used by the indexer to index data from
	// the chain.
	modules []modules.Module
}

func NewIndexer(
	cfg *types.IndexerConfig,
	log log.Logger,
	db database.Database,
	node node.Node,
	modules []modules.Module,
) Indexer {
	return Indexer{
		cfg:          cfg,
		log:          log.With().Str("indexer", cfg.Name).Str("chain_id", node.GetChainID()).Logger(),
		db:           db,
		node:         node,
		HeightsQueue: NewQueue[indexHeight](cfg.HeightQueueSize),
		modules:      modules,
	}
}

// GetName get the name that identify the indexer.
func (i *Indexer) GetName() string {
	return i.cfg.Name
}

// Start starts the indexer.
func (i *Indexer) Start(ctx context.Context, wg *sync.WaitGroup) error {
	// Fetch the current node height
	currentHeigh, err := i.node.GetCurrentHeight(ctx)
	if err != nil {
		return fmt.Errorf("get current node height: %w", err)
	}
	i.log.Debug().Uint64("height", uint64(currentHeigh)).Msg("got current height")

	var startHeight types.Height
	if i.cfg.StartHeight != nil {
		// Start from the height specified in the config
		startHeight = *i.cfg.StartHeight
	} else {
		startHeight = currentHeigh
	}

	// Get the lowest indexed block
	lowestAvailableBlock, err := i.db.GetLowestBlock(i.node.GetChainID())
	if err != nil {
		return fmt.Errorf("get lowest block %w", err)
	}

	// If this is available, we look for missing block from this height.
	if lowestAvailableBlock == nil {
		lowestAvailableBlock = &startHeight
	}

	// Get the blocks that are missing and we need to index
	missingBlocks, err := i.db.GetMissingBlocks(i.node.GetChainID(), *lowestAvailableBlock, currentHeigh)
	if err != nil {
		return fmt.Errorf("get missing blocks: %w", err)
	}
	i.log.Debug().Int("missing", len(missingBlocks)).Msg("got missing blocks")

	// Start the worker that listen for new block produced by the node
	wg.Add(1)
	go i.observeProducedBlocksLoop(ctx, wg, currentHeigh+1, missingBlocks)

	// Starts the indexing workers
	for index := int64(0); index < int64(i.cfg.Workers); index++ {
		wg.Add(1)
		go i.indexingLoop(ctx, wg)
	}

	return nil
}

// FetchAndProcessBlock fetches the block at the provided height and, if fetched successfully, processes it.
func (i *Indexer) FetchAndProcessBlock(ctx context.Context, height types.Height) error {
	i.log.Debug().Uint64("height", uint64(height)).Msg("fetch block")

	// Get the block from the node
	block, err := i.node.GetBlock(ctx, height)
	if err != nil {
		return fmt.Errorf("fetch block %d, %w", height, err)
	}

	// Process the fetched block
	err = i.processBlock(i.log, ctx, block)
	if err != nil {
		return fmt.Errorf("process block %d", height)
	}

	// Save in the database that we have successfully indexed the block
	err = i.db.SaveIndexedBlock(i.node.GetChainID(), height, block.GetTimeStamp())
	if err != nil {
		return fmt.Errorf("save block %d as indexed", height)
	}

	i.log.Debug().Uint64("height", uint64(height)).Msg("block indexed")

	return nil
}

func (i *Indexer) observeProducedBlocksLoop(
	ctx context.Context,
	wg *sync.WaitGroup,
	startHeight types.Height,
	missingHeights []types.Height,
) {
	defer func() {
		i.HeightsQueue.Close()
		i.log.Info().Str("chain-id", i.node.GetChainID()).Msg("stopping node monitoring loop")
		wg.Done()
	}()

	i.log.Info().Str("chain-id", i.node.GetChainID()).Msg("starting node monitoring logic")

	// Enqueue the missing height in the height to index
	for _, height := range missingHeights {
		if !i.HeightsQueue.EnqueueWithContext(ctx, newIndexHeight(height)) {
			break
		}
	}

	// Start the indexing loop
	indexerHeight := startHeight
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !utils.SleepContext(ctx, i.cfg.NodePollingInterval) {
				continue
			}
			// Get the block from the node
			currentNodeHeight, err := i.node.GetCurrentHeight(ctx)
			if err != nil {
				i.log.Err(err).Msg("get current node height")
				continue
			}
			if currentNodeHeight < indexerHeight {
				i.log.Warn().
					Uint64("current", uint64(currentNodeHeight)).
					Uint64("indexer", uint64(indexerHeight)).
					Msg("got node height lower then current indexer height, maybe the node is behind a load balancer")
				continue
			}
			for height := indexerHeight; height <= currentNodeHeight; height++ {
				i.HeightsQueue.EnqueueWithContext(ctx, newIndexHeight(height))
			}
			indexerHeight = currentNodeHeight
		}
	}
}

func (i *Indexer) indexingLoop(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		i.log.Info().Msg("stopping indexing loop")
	}()
	i.log.Info().Msg("started worker")

	for {
		select {
		case <-ctx.Done():
			// exit properly on cancellation
			return
		default:
			indexHeight, ok := i.HeightsQueue.ContextDequeue(ctx)
			if !ok {
				i.log.Warn().Msg("height queue closed, stopping worker")
				return
			}

			// Get the block from the node
			err := i.FetchAndProcessBlock(ctx, indexHeight.Height)
			if err != nil {
				i.log.Err(err).Uint64("height", uint64(indexHeight.Height)).Msg("get and process block")
				i.reEnqueueBlock(ctx, indexHeight)
			}
		}
	}
}

func (i *Indexer) processBlock(_ log.Logger, ctx context.Context, b types.Block) error {
	for _, m := range i.modules {
		// Run the block handling logic
		if blockHandler, ok := m.(modules.BlockHandleModule); ok {
			err := blockHandler.HandleBlock(ctx, b)
			if err != nil {
				return fmt.Errorf("handle block, module: %s err: %w", m.GetName(), err)
			}
		}

		// Run the tx handling logic
		if txHandler, ok := m.(modules.TxHandleModule); ok {
			for _, tx := range b.GetTxs() {
				err := txHandler.HandleTx(ctx, b, tx)
				if err != nil {
					return fmt.Errorf("handle tx, module: %s, tx: %s err: %w", m.GetName(), tx.GetHash(), err)
				}
			}
		}
	}

	return nil
}

func (i *Indexer) reEnqueueBlock(ctx context.Context, indexHeight indexHeight) {
	select {
	case <-ctx.Done():
		i.log.Debug().Uint64("height", uint64(indexHeight.Height)).Msg("skip re-enqueue, context canceled")
	default:
		indexHeight.Attempts += 1
		if indexHeight.Attempts >= i.cfg.MaxAttempts {
			i.log.Error().Uint64("height", uint64(indexHeight.Height)).Msg("failed to parse block, reached max attempts")
			return
		}

		i.log.Info().Uint64("height", uint64(indexHeight.Height)).Msg("re-enqueue block")
		i.HeightsQueue.DelayedEnqueue(ctx, time.Duration(i.cfg.TimeBeforeRetry), indexHeight)
	}
}
