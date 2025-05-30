package rpc

import (
	"time"

	cosmostypes "github.com/milkyway-labs/flux/cosmos/types"
	"github.com/milkyway-labs/flux/types"
)

type StatusRequest struct{}

type NodeInfo struct {
	Network string `json:"network"`
}

type SyncInfo struct {
	LatestBlockHeight   types.Height `json:"latest_block_height,string"`
	LatestBlockTime     time.Time    `json:"latest_block_time"`
	EarliestBlockHeight types.Height `json:"earliest_block_height,string"`
	EarliestBlockTime   time.Time    `json:"earliest_block_time"`
}

type StatusResponse struct {
	NodeInfo NodeInfo `json:"node_info"`
	SyncInfo SyncInfo `json:"sync_info"`
}

type BlockRequest struct {
	Height *types.Height `json:"height,string,omitempty"`
}

type BlockResponse struct {
	Block Block `json:"block"`
}

type BlockHeader struct {
	ChainID string       `json:"chain_id"`
	Height  types.Height `json:"height,string"`
	Time    time.Time    `json:"time"`
}

type Block struct {
	BlockHeader `json:"header"`
}

type BlockResultsRequest struct {
	Height *types.Height `json:"height,string,omitempty"`
}

type BlockResultsResponse struct {
	Height              types.Height           `json:"height,string"`
	TxsResults          []ResponseDeliverTx    `json:"txs_results"`
	BeginBlockEvents    cosmostypes.ABCIEvents `json:"begin_block_events"`
	EndBlockEvents      cosmostypes.ABCIEvents `json:"end_block_events"`
	FinalizeBlockEvents cosmostypes.ABCIEvents `json:"finalize_block_events"`
}

type ResponseDeliverTx struct {
	Code      uint32                 `json:"code"`
	Data      []byte                 `json:"data"`
	TxHash    string                 `json:"txhash"`
	Log       string                 `json:"log"`
	GasWanted int64                  `json:"gas_wanted,string"`
	GasUsed   int64                  `json:"gas_used,string"`
	Events    cosmostypes.ABCIEvents `json:"events"`
}

func (resp ResponseDeliverTx) IsOK() bool {
	return resp.Code == 0
}
