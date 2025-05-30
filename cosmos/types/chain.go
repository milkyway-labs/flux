package types

import (
	"time"

	"github.com/milkyway-labs/flux/types"
)

// ----------------------------------------------------------------------------
// -- Block related data structures
// ----------------------------------------------------------------------------

type BlockHeader struct {
	ChainID string
	Height  types.Height
	Time    time.Time
}

func NewBlockHeader(chainID string, height types.Height, time time.Time) BlockHeader {
	return BlockHeader{
		ChainID: chainID,
		Height:  height,
		Time:    time,
	}
}

type Block struct {
	Header              BlockHeader
	Txs                 []Tx
	BeginBlockEvents    ABCIEvents
	EndBlockEvents      ABCIEvents
	FinalizeBlockEvents ABCIEvents
}

var _ types.Block = &Block{}

func NewBlock(
	header BlockHeader,
	txs []Tx,
	beginBlockEvents ABCIEvents,
	endBlockEvents ABCIEvents,
	finalizeBlockEvents ABCIEvents,
) *Block {
	return &Block{
		Header:              header,
		Txs:                 txs,
		BeginBlockEvents:    beginBlockEvents,
		EndBlockEvents:      endBlockEvents,
		FinalizeBlockEvents: finalizeBlockEvents,
	}
}

// GetChainID implements types.Block.
func (b *Block) GetChainID() string {
	return b.Header.ChainID
}

// GetHeight implements types.Block.
func (b *Block) GetHeight() types.Height {
	return b.Header.Height
}

// GetTimeStamp implements types.Block.
func (b *Block) GetTimeStamp() time.Time {
	return b.Header.Time
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
// -- Tx related data structures
// ----------------------------------------------------------------------------

var _ types.Tx = &Tx{}

type Tx struct {
	Code   uint32
	Data   []byte
	TxHash string
	Events ABCIEvents
	Log    string
}

func NewTx(
	code uint32,
	data []byte,
	hash string,
	events ABCIEvents,
	log string,
) Tx {
	return Tx{
		Code:   code,
		Data:   data,
		TxHash: hash,
		Events: events,
		Log:    log,
	}
}

// GetHash implements types.Tx.
func (t *Tx) GetHash() string {
	return t.TxHash
}

// IsSuccessful implements types.Tx.
func (t *Tx) IsSuccessful() bool {
	return t.Code == 0
}
