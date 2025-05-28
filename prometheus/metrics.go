package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

// WorkersCount represents the Telemetry counter used to track the number of active workers
var WorkersCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "indexer_active_workers",
		Help: "Number of active workers for each indexer.",
	},
	[]string{"indexer_name"},
)

// LatestIndexedHeightByIndexer represents the Telemetry counter used to track
// the last indexed height for each indexer.
var LatestIndexedHeightByIndexer = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "indexer_latest_indexed_height",
		Help: "Height of the last indexed block.",
	},
	[]string{"indexer_name"},
)

// IndexerFailedBlocks represents the Telemetry counter used to track the failed blocks
var IndexerFailedBlocks = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "indexer_failed_blocks",
		Help: "Height of the block for which the indexer has failed.",
	},
	[]string{"indexer_name"},
)

func init() {
	prometheus.MustRegister(WorkersCount)
	prometheus.MustRegister(LatestIndexedHeightByIndexer)
	prometheus.MustRegister(IndexerFailedBlocks)
}
