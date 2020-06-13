package eth

import (
	"cos-backend-com/src/common/flake"
)

type TransactionState int

const (
	TransactionStateNull        TransactionState = 0
	TransactionStateWaitConfirm TransactionState = 1
	TransactionStateSuccess     TransactionState = 2
	TransactionStateFailed      TransactionState = 3
)

func (t TransactionState) Validate() bool {
	switch t {
	case TransactionStateNull, TransactionStateWaitConfirm, TransactionStateSuccess, TransactionStateFailed:
		return true
	}
	return false
}

type TransactionSource string

const (
	TransactionSourceStartup        TransactionSource = "startup"
	TransactionSourceStartupSetting TransactionSource = "startupSetting"
)

func (t TransactionSource) Validate() bool {
	switch t {
	case TransactionSourceStartup, TransactionSourceStartupSetting:
		return true
	}
	return false
}

type CreateTransactionsInput struct {
	TxId     string            `json:"txId"`
	Source   TransactionSource `json:"source"`
	SourceId flake.ID          `json:"sourceId"`
}

type UpdateTransactionsInput struct {
	BlockAddr string `json:"blockAddr"`
	State     int    `json:"state"`
}

type TransactionsResult struct {
	TxId      string `json:"txId"`
	BlockAddr string `json:"blockAddr"`
}
