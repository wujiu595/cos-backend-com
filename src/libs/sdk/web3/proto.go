package web3

type EcrecoverInput struct {
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}

type EcrecoverOutput struct {
	PublicKey string `json:"publicKey"`
}
