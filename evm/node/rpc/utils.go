package rpc

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	evmtypes "github.com/milkyway-labs/flux/evm/types"
	"github.com/milkyway-labs/flux/utils"
)

// Convert a go-ethereum log entry into a log entry with our format.
func LogFromEthlog(log *ethtypes.Log) evmtypes.LogEntry {
	return evmtypes.LogEntry{
		Address: log.Address.Bytes(),
		Topics: utils.Map(log.Topics, func(hash common.Hash) evmtypes.EVMBytes {
			return hash.Bytes()
		}),
		Data:        log.Data,
		BlockNumber: log.BlockNumber,
		TxIndex:     log.TxIndex,
		Index:       log.Index,
		Removed:     log.Removed,
	}
}

func TxFromReceipt(r *ethtypes.Receipt) *evmtypes.Tx {
	return evmtypes.NewTx(r.TxHash.Bytes(), r.Status == ethtypes.ReceiptStatusSuccessful)
}
