package indexer

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/milkyway-labs/flux/node"
	"github.com/milkyway-labs/flux/types"
	"github.com/milkyway-labs/flux/utils"
)

// IndexerHeight represents the height structure used by the
// indexer to know which height to fetch and monitor its attempts.
type IndexerHeight struct {
	Height   types.Height
	Attempts uint32
}

func NewIndexerHeight(height types.Height) IndexerHeight {
	return IndexerHeight{
		Height:   height,
		Attempts: 0,
	}
}

// HeightProducer represents a generic component capable of
// providing heights to be fetched and parsed by a Worker.
type HeightProducer interface {
	// EnqueueHeights is called when the component needs to
	// provide the heights to parse. The heights should be enqueued
	// into the provided Queue struct.
	EnqueueHeights(ctx context.Context, queue *Queue[IndexerHeight]) error
}

// ----------------------------------------------------------------------------
// ---- Combined producer height producer
// ----------------------------------------------------------------------------

var _ HeightProducer = &CombinedHeightProducer{}

// CombinedHeightProducer represents a HeightProducer that combines multiple HeightProducers
// into a single one. It calls the provided producers from the first to the last,
// following the order in which they were added by the user.
type CombinedHeightProducer struct {
	producers []HeightProducer
}

// NewCombinedHeightProducer creates a new CombinedHeightProducer instance.
func NewCombinedHeightProducer(producers ...HeightProducer) *CombinedHeightProducer {
	return &CombinedHeightProducer{
		producers: producers,
	}
}

// AddProducer adds a producer to the list of producers to call.
func (c *CombinedHeightProducer) AddProducer(producer HeightProducer) *CombinedHeightProducer {
	c.producers = append(c.producers, producer)
	return c
}

func (c *CombinedHeightProducer) EnqueueHeights(ctx context.Context, queue *Queue[IndexerHeight]) error {
	for _, p := range c.producers {
		err := p.EnqueueHeights(ctx, queue)
		if err != nil {
			return err
		}
	}

	return nil
}

// ----------------------------------------------------------------------------
// ---- Range height producer
// ----------------------------------------------------------------------------

var _ HeightProducer = &RangeHeightProducer{}

// RangeHeightProducer is a HeightProducer that produces all the heights
// from the specified `from` height to the `to` height, inclusive.
type RangeHeightProducer struct {
	from types.Height
	to   types.Height
}

// NewRangeHeightProducer creates a new HeightProducer instance.
func NewRangeHeightProducer(from types.Height, to types.Height) *RangeHeightProducer {
	return &RangeHeightProducer{
		from: from,
		to:   to,
	}
}

func (r *RangeHeightProducer) EnqueueHeights(ctx context.Context, queue *Queue[IndexerHeight]) error {
	for i := r.from; i <= r.to; i++ {
		if !queue.EnqueueWithContext(ctx, NewIndexerHeight(i)) {
			break
		}
	}

	return nil
}

// ----------------------------------------------------------------------------
// ---- List height producer
// ----------------------------------------------------------------------------

var _ HeightProducer = &ListHeightProducer{}

// ListHeightProducer is a HeightProducer that produces the heights from a list
// provided by the user.
type ListHeightProducer struct {
	heights []types.Height
}

// NewListHeightProducer creates a new ListHeightProducer instance.
func NewListHeightProducer(heights []types.Height) *ListHeightProducer {
	return &ListHeightProducer{
		heights: heights,
	}
}

func (r *ListHeightProducer) EnqueueHeights(ctx context.Context, queue *Queue[IndexerHeight]) error {
	for _, h := range r.heights {
		if !queue.EnqueueWithContext(ctx, NewIndexerHeight(h)) {
			break
		}
	}

	return nil
}

// ----------------------------------------------------------------------------
// ---- NodeHeightProducer
// ----------------------------------------------------------------------------

var _ HeightProducer = &NodeHeightProducer{}

// NodeHeightProducer is a HeightProducer that produces heights by monitoring
// newly produced blocks by a node after the configured height.
// The node is monitored by querying its latest available block at the user-provided polling interval.
type NodeHeightProducer struct {
	logger          zerolog.Logger
	from            types.Height
	pollingInterval time.Duration
	node            node.Node
}

// NewNodeHeightProducer creates a new NodeHeightProducer instance.
func NewNodeHeightProducer(
	logger zerolog.Logger,
	node node.Node,
	pollingInterval time.Duration,
	from types.Height,
) *NodeHeightProducer {
	logger = logger.With().
		Str("component", "NodeHeightProducer").
		Str("chain-id", node.GetChainID()).
		Logger()

	return &NodeHeightProducer{
		logger:          logger,
		node:            node,
		pollingInterval: pollingInterval,
		from:            from,
	}
}

func (n *NodeHeightProducer) EnqueueHeights(ctx context.Context, queue *Queue[IndexerHeight]) error {
	defer func() {
		n.logger.Info().Msg("stopping node monitoring loop")
	}()

	toFetchHeight := n.from
	n.logger.Info().
		Uint64("start height", uint64(toFetchHeight)).
		Msg("start node monitoring loop")
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if !utils.SleepContext(ctx, n.pollingInterval) {
				continue
			}
			// Get the block from the node
			currentNodeHeight, err := n.node.GetCurrentHeight(ctx)
			if err != nil {
				n.logger.Err(err).Msg("get current node height")
				continue
			}
			// Prevent enqueue of an already indexed block
			if toFetchHeight > currentNodeHeight {
				continue
			}
			for ; toFetchHeight <= currentNodeHeight; toFetchHeight++ {
				queue.EnqueueWithContext(ctx, NewIndexerHeight(toFetchHeight))
			}
			// After the loop, toFetchHeight is currentNodeHeight + 1 now
		}
	}
}
