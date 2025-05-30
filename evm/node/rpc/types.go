package rpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/milkyway-labs/flux/evm/types"
)

type GetBlockBlockByNumberResponse struct {
	Difficulty       *hexutil.Big    `json:"difficulty"`
	ExtraData        hexutil.Bytes   `json:"extraData"`
	GasLimit         hexutil.Uint64  `json:"gasLimit"`
	GasUsed          hexutil.Uint64  `json:"gasUsed"`
	Hash             common.Hash     `json:"hash"`
	LogsBloom        hexutil.Bytes   `json:"logsBloom"`
	Miner            common.Address  `json:"miner"`
	MixHash          common.Hash     `json:"mixHash"`
	Number           *hexutil.Uint64 `json:"number"`
	ParentHash       common.Hash     `json:"parentHash"`
	ReceiptsRoot     common.Hash     `json:"receiptsRoot"`
	Sha3Uncles       common.Hash     `json:"sha3Uncles"`
	Size             hexutil.Uint64  `json:"size"`
	StateRoot        common.Hash     `json:"stateRoot"`
	Timestamp        hexutil.Uint64  `json:"timestamp"`
	TotalDifficulty  *hexutil.Big    `json:"totalDifficulty"`
	Transactions     []common.Hash   `json:"transactions"`
	TransactionsRoot common.Hash     `json:"transactionsRoot"`
	Uncles           []common.Hash   `json:"uncles"`
}

type GetLogsRequest struct {
	FromBlock hexutil.Uint64
	ToBlock   hexutil.Uint64
}

type GetLogsResponse []evmtypes.LogEntry
