package indexer

import (
	"context"
	"fmt"
	"sync"

	log "github.com/rs/zerolog"

	"github.com/milkyway-labs/chain-indexer/database"
	"github.com/milkyway-labs/chain-indexer/modules"
	"github.com/milkyway-labs/chain-indexer/node"
	"github.com/milkyway-labs/chain-indexer/types"
)

// Worker represents a component that is responsible of fetching a block
// and pass it to a list of modules to be processed.
type Worker struct {
	// Logger used by the worker.
	log log.Logger
	// Config of the indexer to which this worker belongs.
	cfg *types.IndexerConfig
	// Queue used by the worker to retrieve the height of the
	// blocks to index.
	heightsQueue *Queue[IndexerHeight]
	// Database used by the indexer to store its state.
	db database.Database
	// The node from which the indexer fetches the blocks to index.
	node node.Node
	// List of modules that will be used by the indexer to index data from
	// the chain.
	modules []modules.Module
}

func NewWorker(
	cfg *types.IndexerConfig,
	log log.Logger,
	heightsQueue *Queue[IndexerHeight],
	db database.Database,
	node node.Node,
	modules []modules.Module,
) Worker {
	return Worker{
		cfg:          cfg,
		log:          log.With().Str("component", "worker").Logger(),
		heightsQueue: heightsQueue,
		db:           db,
		node:         node,
		modules:      modules,
	}
}

// Start the worker logic.
func (i *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go i.workerLoop(ctx, wg)
}

func (i *Worker) workerLoop(ctx context.Context, wg *sync.WaitGroup) {
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
			indexHeight, ok := i.heightsQueue.ContextDequeue(ctx)
			if !ok {
				i.log.Warn().Msg("height queue closed, stopping worker")
				return
			}

			// Get the block from the node
			err := i.fetchAndProcessBlock(ctx, indexHeight.Height)
			if err != nil {
				i.log.Err(err).Uint64("height", uint64(indexHeight.Height)).Msg("get and process block")
				i.reEnqueueBlock(ctx, indexHeight)
			}
		}
	}
}

// fetchAndProcessBlock fetches the block at the provided height and, if fetched successfully, processes it.
func (i *Worker) fetchAndProcessBlock(ctx context.Context, height types.Height) error {
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

func (i *Worker) processBlock(_ log.Logger, ctx context.Context, b types.Block) error {
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

func (i *Worker) reEnqueueBlock(ctx context.Context, indexHeight IndexerHeight) {
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
		i.heightsQueue.DelayedEnqueue(ctx, i.cfg.TimeBeforeRetry, indexHeight)
	}
}
