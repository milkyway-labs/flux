package types

import (
	"encoding/hex"
	"math/big"
	"slices"
)

type EVMBytes []byte

// TrimLeadingZero removes the leading zeros from the bytes.
func (e EVMBytes) TrimLeadingZero() EVMBytes {
	// Find the first non-zero byte
	firstNonZero := slices.IndexFunc(e, func(b byte) bool {
		return b != 0
	})

	// Handle all-zero case
	if firstNonZero == -1 {
		return []byte{0}
	}

	// Remove the leading zeros
	return e[firstNonZero:]
}

// Hex hex encode the bytes with a 0x prefix.
func (e EVMBytes) Hex() string {
	return "0x" + hex.EncodeToString(e)
}

// NormalizedHex hex encoded the bytes without the leading zeros.
func (e EVMBytes) NormalizedHex() string {
	return e.TrimLeadingZero().Hex()
}

// Int convert the bytes to a *big.Int.
func (e EVMBytes) Int() *big.Int {
	return new(big.Int).SetBytes(e)
}
