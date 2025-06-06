package rpc

import "crypto/sha256"

// TxHasher represents a function that given the bytes of a transaction
// returns its hash.
type TxHasher func(txData []byte) []byte

// DefaultTxHasher is the default TxHasher used by the cosmos rpc node
func DefaultTxHasher(txData []byte) []byte {
	hash := sha256.Sum256(txData)
	return hash[:]
}
