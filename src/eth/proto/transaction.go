package proto

import (
	"cos-backend-com/src/common/flake"
	ethSdk "cos-backend-com/src/libs/sdk/eth"
)

type TransactionsOutput struct {
	Id        flake.ID                 `json:"id" db:"id"`                // id (PK)
	TxId      string                   `json:"txId" db:"tx_id"`           // tx_id
	BlockAddr string                   `json:"blockAddr" db:"block_addr"` // block_addr
	Source    ethSdk.TransactionSource `json:"source" db:"source"`        // source
	SourceId  flake.ID                 `json:"sourceId" db:"source_id"`   // source_id
	RetryTime int                      `json:"retryTime" db:"retry_time"` // retry_time
	State     ethSdk.TransactionState  `json:"state" db:"state"`          // state
}
