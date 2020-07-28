package processor

import (
	"context"
	"cos-backend-com/src/eth/proto"
	"cos-backend-com/src/libs/models/ethmodels"
	ethSdk "cos-backend-com/src/libs/sdk/eth"

	"github.com/qiniu/x/log"
)

type Updater struct {
	TransactionInput <-chan *proto.TransactionsOutput
}

func (c *Updater) Process() {
	for transactionInput := range c.TransactionInput {
		if transactionInput.State == ethSdk.TransactionStateSuccess {
			updateTransactionsInput := ethSdk.UpdateTransactionsInput{
				BlockAddr: transactionInput.BlockAddr,
				State:     transactionInput.State,
			}
			if err := ethmodels.Transactions.UpdateWithConfirmSource(context.Background(), transactionInput.Id, transactionInput.SourceId, transactionInput.Source, &updateTransactionsInput); err != nil {
				log.Warn(err)
			}
		} else {
			updateTransactionsInput := ethSdk.UpdateTransactionsInput{
				BlockAddr: transactionInput.BlockAddr,
				State:     transactionInput.State,
			}
			if err := ethmodels.Transactions.Update(context.Background(), transactionInput.Id, &updateTransactionsInput); err != nil {
				log.Warn(err)
			}
		}
	}
}
