package processor

import (
	"cos-backend-com/src/eth"

	"github.com/ethereum/go-ethereum/ethclient"
)

var EthClient *ethclient.Client

func InitEthClient() {
	ethClient, err := ethclient.Dial(eth.Env.EthClient.EndPoint + "/" + eth.Env.EthClient.InfuraKey)
	if err != nil {
		panic(err)
	}
	EthClient = ethClient
}
