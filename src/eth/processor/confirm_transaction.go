package processor

import (
	"context"
	"cos-backend-com/src/eth/proto"
	ethSdk "cos-backend-com/src/libs/sdk/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type Confirmer struct {
	TransactionInput  <-chan *proto.TransactionsOutput
	TransactionOutput chan<- *proto.TransactionsOutput
}

func (c *Confirmer) Process() {
	for transactionInput := range c.TransactionInput {
		txHash := common.HexToHash(transactionInput.TxId)
		receipt, err := EthClient.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			transactionInput.State = ethSdk.TransactionStateFailed
			log.Warn(err.Error())
		} else {
			transactionInput.BlockAddr = receipt.BlockHash.Hex()
			transactionInput.State = ethSdk.TransactionStateSuccess
		}
		c.TransactionOutput <- transactionInput
	}
}
