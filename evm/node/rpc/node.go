package rpc

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/rs/zerolog"

	evmtypes "github.com/milkyway-labs/chain-indexer/evm/types"
	"github.com/milkyway-labs/chain-indexer/node"
	"github.com/milkyway-labs/chain-indexer/rpc/jsonrpc2"
	"github.com/milkyway-labs/chain-indexer/types"
)

var _ node.Node = &Node{}

// Node represents a node instance capable of fetching data from an EVM node
// trough JSON-RPC calls.
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

	var chainIDHex string
	if err := jsonRPCClient.Call(ctx, "eth_chainId", []string{}, &chainIDHex); err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}

	return &Node{
		cfg:     cfg,
		logger:  logger.With().Str(NodeType, cfg.URL).Logger(),
		client:  jsonRPCClient,
		chainID: chainIDHex,
	}, nil
}

// GetBlock implements node.Node.
func (n *Node) GetBlock(ctx context.Context, height types.Height) (types.Block, error) {
	ethBlock, err := n.GetEthBlock(ctx, height)
	if err != nil {
		return nil, fmt.Errorf("get eth block failed: %w", err)
	}

	// Fetch the log for this block
	logs, err := n.GetLogs(ctx, height)
	if err != nil {
		return nil, fmt.Errorf("get eth logs")
	}

	// Create the Tx objects from the logs result
	txs := make(map[common.Hash]evmtypes.Tx)
	for _, logEntry := range logs {
		tx, ok := txs[logEntry.TransactionHash]
		// Create a new Tx object
		if !ok {
			tx = evmtypes.NewTx(logEntry.TransactionHash.Hex())
		}
		// Add the relevant log entry to the tx logs
		tx.Logs = append(tx.Logs, logEntry)

		txs[logEntry.TransactionHash] = tx
	}

	// Create the block
	block := evmtypes.NewBlock(
		n.chainID,
		height,
		time.Unix(int64(ethBlock.Timestamp), 0),
		evmtypes.Logs(logs),
		slices.Collect(maps.Values(txs)),
	)

	return &block, nil
}

// GetChainID implements node.Node.
func (n *Node) GetChainID() string {
	return n.chainID
}

// GetCurrentHeight implements node.Node.
func (n *Node) GetCurrentHeight(ctx context.Context) (types.Height, error) {
	var blockNumberHex string
	if err := n.client.Call(ctx, "eth_blockNumber", []string{}, &blockNumberHex); err != nil {
		return 0, fmt.Errorf("get current height: %w", err)
	}

	// Remove the 0x prefix and parse as hex number
	height, err := strconv.ParseUint(blockNumberHex[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("parse hex height: %w", err)
	}

	return types.Height(height), nil
}

// GetLowestHeight implements node.Node.
func (n *Node) GetLowestHeight(ctx context.Context) (types.Height, error) {
	// Try to get the first block
	_, err := n.GetEthBlock(ctx, 0)
	if err == nil {
		return 0, nil
	}

	// Get the current height
	currentHeight, err := n.GetCurrentHeight(ctx)
	if err != nil {
		return 0, err
	}

	// Start a binary search for the lowest available block
	low := types.Height(0)
	high := currentHeight
	lowestAvailable := currentHeight

	for low <= high {
		mid := (low + high) / 2
		if _, err := n.GetEthBlock(ctx, mid); err == nil {
			lowestAvailable = mid
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	return lowestAvailable, nil
}

// Performs a "eth_getBlockByNumber" with the provided height.
func (n *Node) GetEthBlock(ctx context.Context, height types.Height) (GetBlockBlockByNumberResponse, error) {
	var response GetBlockBlockByNumberResponse
	if err := n.client.Call(ctx, "eth_getBlockByNumber", []any{hexutil.Uint64(height), false}, &response); err != nil {
		return GetBlockBlockByNumberResponse{}, fmt.Errorf("get current height: %w", err)
	}

	return response, nil
}

// Performs a "eth_getLogs" with the provided height.
func (n *Node) GetLogs(ctx context.Context, height types.Height) (GetLogsResponse, error) {
	var response GetLogsResponse
	if err := n.client.Call(ctx, "eth_getLogs", []any{GetLogsRequest{
		FromBlock: hexutil.Uint64(height),
		ToBlock:   hexutil.Uint64(height),
	}}, &response); err != nil {
		return nil, fmt.Errorf("get logs: %w", err)
	}

	return response, nil
}
