package web3

import "net/http"

func NewWeb3Service(host string) interface{} {
	return func(adminRt http.RoundTripper) Web3Service {
		return New(Config{
			Host:      host,
			Transport: adminRt,
		})
	}
}
