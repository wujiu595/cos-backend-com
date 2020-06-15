package handlers

import (
	"context"
	"cos-backend-com/src/eth/proto"
	"cos-backend-com/src/libs/models/ethmodels"

	s "github.com/wujiu2020/strip"
)

func LoadPayload(ctx context.Context, strip *s.Strip, blockAddr string, in chan<- *proto.TransactionsOutput) error {
	var inputs []*proto.TransactionsOutput
	if err := ethmodels.Transactions.List(ctx, &inputs); err != nil {
		strip.Logger().Error(err)
		return err
	}
	for _, input := range inputs {
		in <- input
	}
	return nil
}
