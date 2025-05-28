package rpc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"

	cosmostypes "github.com/milkyway-labs/chain-indexer/cosmos/types"
	"github.com/milkyway-labs/chain-indexer/node"
	"github.com/milkyway-labs/chain-indexer/rpc/jsonrpc2"
	"github.com/milkyway-labs/chain-indexer/types"
)

var _ node.Node = &Node{}

type Node struct {
	cfg     Config
	logger  zerolog.Logger
	client  *jsonrpc2.Client
	chainID string
}

func NewNode(ctx context.Context, logger zerolog.Logger, cfg Config) (*Node, error) {
	jsonRPCClient, err := jsonrpc2.NewClient(cfg.URL, &http.Client{
		Timeout: cfg.RequestTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("create rpc client: %w", err)
	}

	var res StatusResponse
	if err := jsonRPCClient.Call(ctx, "status", StatusRequest{}, &res); err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}

	return &Node{
		cfg:     cfg,
		logger:  logger.With().Str("cosmos-node", cfg.URL).Logger(),
		client:  jsonRPCClient,
		chainID: res.NodeInfo.Network,
	}, nil
}

// GetChainID implements node.Node.
func (r *Node) GetChainID() string {
	return r.chainID
}

// GetCurrentHeight implements node.Node.
func (r *Node) GetCurrentHeight(ctx context.Context) (types.Height, error) {
	var res StatusResponse
	if err := r.client.Call(ctx, "status", StatusRequest{}, &res); err != nil {
		return 0, fmt.Errorf("call status: %w", err)
	}

	return res.SyncInfo.LatestBlockHeight, nil
}

// GetLowestHeight implements node.Node.
func (r *Node) GetLowestHeight(ctx context.Context) (types.Height, error) {
	var res StatusResponse
	if err := r.client.Call(ctx, "status", StatusRequest{}, &res); err != nil {
		return 0, fmt.Errorf("call status: %w", err)
	}

	return res.SyncInfo.EarliestBlockHeight, nil
}

// GetBlock implements node.Node.
func (r *Node) GetBlock(ctx context.Context, height types.Height) (types.Block, error) {
	var blockResponse BlockResponse
	if err := r.client.Call(ctx, "block", BlockRequest{Height: &height}, &blockResponse); err != nil {
		return nil, fmt.Errorf("call block: %w", err)
	}

	var blockResultsResponse BlockResultsResponse
	if err := r.client.Call(ctx, "block_results", BlockResultsRequest{Height: &height}, &blockResultsResponse); err != nil {
		return nil, fmt.Errorf("call block_results: %w", err)
	}

	// Extract the tx events
	txs := make([]cosmostypes.Tx, len(blockResultsResponse.TxsResults))
	for txIndex, txResult := range blockResultsResponse.TxsResults {
		var txEvents cosmostypes.ABCIEvents
		if r.cfg.TxEventsFromLog {
			// We should parse the events from the log, ensure the transaction
			// was successful before parsing the log
			if txResult.Code == 0 {
				parsedEvents, err := ParseEventsFromTxLog(txResult.Log)
				if err != nil {
					return nil, fmt.Errorf("parse tx.log (height %d, txIndex %d): %w", height, txIndex, err)
				}
				txEvents = parsedEvents
			}
		} else {
			txEvents = txResult.Events
		}

		txs[txIndex] = cosmostypes.NewTx(
			txResult.Code,
			txResult.Data,
			txResult.TxHash,
			txEvents,
			txResult.Log,
		)
	}

	// Decode the block events attributes
	if r.cfg.DecodeBlockEventAttributes {
		decoded, err := DecodeABCIEvents(blockResultsResponse.BeginBlockEvents)
		if err != nil {
			return nil, fmt.Errorf("decode begin block events (height: %d): %w", height, err)
		}
		blockResultsResponse.BeginBlockEvents = decoded

		decoded, err = DecodeABCIEvents(blockResultsResponse.EndBlockEvents)
		if err != nil {
			return nil, fmt.Errorf("decode end block events (height: %d): %w", height, err)
		}
		blockResultsResponse.EndBlockEvents = decoded

		decoded, err = DecodeABCIEvents(blockResultsResponse.FinalizeBlockEvents)
		if err != nil {
			return nil, fmt.Errorf("decode finalize block events (height: %d): %w", height, err)
		}
		blockResultsResponse.FinalizeBlockEvents = decoded
	}

	if len(blockResultsResponse.FinalizeBlockEvents) > 0 {
		// In case we have the finalized blocks let's extract the begin and end
		// block events.
		var beginBlocksEvents cosmostypes.ABCIEvents
		var endBlockEvents cosmostypes.ABCIEvents
		for _, event := range blockResultsResponse.FinalizeBlockEvents {
			// Check if the event is a BeginBlock event
			if _, found := event.FindAttributeFunc(func(a cosmostypes.ABCIEventAttribute) bool {
				return a.Key == "mode" && a.Value == "BeginBlock"
			}); found {
				beginBlocksEvents = append(beginBlocksEvents, event)
			}

			// Check if the event is an EndBlock event
			if _, found := event.FindAttributeFunc(func(a cosmostypes.ABCIEventAttribute) bool {
				return a.Key == "mode" && a.Value == "EndBlock"
			}); found {
				endBlockEvents = append(endBlockEvents, event)
			}
		}

		// Add the events to the begin and end block events
		blockResultsResponse.BeginBlockEvents = append(blockResultsResponse.BeginBlockEvents, beginBlocksEvents...)
		blockResultsResponse.EndBlockEvents = append(blockResultsResponse.EndBlockEvents, endBlockEvents...)
	}

	blockHeader := cosmostypes.NewBlockHeader(blockResponse.Block.ChainID, blockResponse.Block.Height, blockResponse.Block.Time)
	return cosmostypes.NewBlock(
		blockHeader,
		txs,
		blockResultsResponse.BeginBlockEvents,
		blockResultsResponse.EndBlockEvents,
		blockResultsResponse.FinalizeBlockEvents,
	), nil
}
