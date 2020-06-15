package main

import (
	"context"
	"cos-backend-com/src/eth/cmd/eth/app"
	"cos-backend-com/src/eth/handlers"
	"cos-backend-com/src/eth/processor"
	"cos-backend-com/src/eth/proto"
	"flag"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/trustmaster/goflow"
	s "github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/utils/helpers"
)

func main() {
	strip := s.New()

	var (
		confPath string
		files    helpers.SliceValue
	)
	flagSet := flag.NewFlagSet("app", flag.ExitOnError)
	flagSet.StringVar(&confPath, "conf", "", "")
	flagSet.Var(&files, "f", "")
	flagSet.Parse(os.Args[1:])

	app.AppInit(strip, confPath, files...).Start()
	go confirmTransaction(strip)
	strip.Run()
}

func confirmTransaction(strip *s.Strip) {
	transactionFlow := processor.NewConfirmingApp()
	in := make(chan *proto.TransactionsOutput, 10)
	transactionFlow.SetInPort("In", in)
	wait := goflow.Run(transactionFlow)
	go func() {
		headers := make(chan *types.Header)
		sub, err := processor.EthClient.SubscribeNewHead(context.Background(), headers)
		if err != nil {
			panic(err)
		}
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case header := <-headers:
				if err := handlers.LoadPayload(context.TODO(), strip, header.Hash().Hex(), in); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()
	<-wait
}
