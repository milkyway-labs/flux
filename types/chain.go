package types

import (
	"math"
	"time"
)

// Height type alias used to represent a block chain height.
type Height uint64

const MaxHeight Height = Height(math.MaxUint64)

// Block represents a generic block produced by a blockchain.
type Block interface {
	// GetChainID provides the ID of the blockchain that produced this Block.
	GetChainID() string
	// GetHeight provides the height at which this block has been produced.
	GetHeight() Height
	// GetTimeStamp provides the time at which this block has been produced.
	GetTimeStamp() time.Time
	// GetTxs get the transactions included in this Block.
	GetTxs() []Tx
}

// Tx represents a generic transaction that has been included into a Block.
type Tx interface {
	// GetHash gets this transaction hash.
	GetHash() string
	// IsSuccessful returns true if the transaction has been executed without errors, false otherwise.
	IsSuccessful() bool
}
