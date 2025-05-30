package types

import (
	"slices"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/milkyway-labs/flux/types"
	"github.com/milkyway-labs/flux/utils"
)

var _ types.Block = &Block{}

type Block struct {
	ChainID   string
	Height    types.Height
	Timestamp time.Time
	Logs      Logs
	Txs       []Tx
}

func NewBlock(
	chainID string,
	height types.Height,
	timestamp time.Time,
	logs Logs,
	txs []Tx,
) Block {
	return Block{
		ChainID:   chainID,
		Height:    height,
		Timestamp: timestamp,
		Logs:      logs,
		Txs:       txs,
	}
}

// GetChainID implements types.Block.
func (b *Block) GetChainID() string {
	return b.ChainID
}

// GetHeight implements types.Block.
func (b *Block) GetHeight() types.Height {
	return b.Height
}

// GetTimeStamp implements types.Block.
func (b *Block) GetTimeStamp() time.Time {
	return b.Timestamp
}

// GetTxs implements types.Block.
func (b *Block) GetTxs() []types.Tx {
	result := make([]types.Tx, len(b.Txs))
	for i, tx := range b.Txs {
		result[i] = &tx
	}
	return result
}

// ----------------------------------------------------------------------------
// ---- EVM Tx type definition
// ----------------------------------------------------------------------------

var _ types.Tx = &Tx{}

type Tx struct {
	Hash string
	Logs Logs
}

func NewTx(hash string) Tx {
	return Tx{
		Hash: hash,
	}
}

// GetHash implements types.Tx.
func (t *Tx) GetHash() string {
	return t.Hash
}

// IsSuccessful implements types.Tx.
func (t *Tx) IsSuccessful() bool {
	return true
}

// ----------------------------------------------------------------------------
// ---- EVM Log entry
// ----------------------------------------------------------------------------

// LogEntry represents an Ethereum event log entry.
type LogEntry struct {
	Address          common.Address `json:"address"`
	Topics           []common.Hash  `json:"topics"`
	Data             hexutil.Bytes  `json:"data"`
	BlockNumber      hexutil.Uint64 `json:"blockNumber"`
	TransactionHash  common.Hash    `json:"transactionHash"`
	TransactionIndex hexutil.Uint64 `json:"transactionIndex"`
	BlockHash        common.Hash    `json:"blockHash"`
	LogIndex         hexutil.Uint64 `json:"logIndex"`
	Removed          bool           `json:"removed"`
}

// ----------------------------------------------------------------------------
// ---- EVM Logs
// ----------------------------------------------------------------------------

type Logs []LogEntry

// FindEntryFunc finds the first LogEntry matching the given predicate.
func (l Logs) FindEntryFunc(predicate func(LogEntry) bool) (LogEntry, bool) {
	index := slices.IndexFunc(l, predicate)
	if index == -1 {
		return LogEntry{}, false
	}
	return l[index], true
}

// FilterFunc returns all the LogEntry that match the given predicate.
func (l Logs) FilterFunc(predicate func(LogEntry) bool) Logs {
	return utils.Filter(l, predicate)
}
