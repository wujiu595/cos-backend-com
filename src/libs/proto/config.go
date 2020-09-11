package proto

type ServiceEndpoint struct {
	Account      string `conf:"account"`
	Cores        string `conf:"cores"`
	Notification string `conf:"notification"`
	Eth          string `conf:"eth"`
	Web3         string `conf:"web3"`
}
