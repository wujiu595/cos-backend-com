package cores

import (
	"cos-backend-com/src/common/flake"
)

type PrepareIdOutput struct {
	Id flake.ID `json:"id"`
}

type Token struct {
	Name   string `json:"name" db:"name"`
	Symbol string `json:"symbol" db:"symbol"`
}

type PayTokenListOutput struct {
	PayTokens []Token `json:"payTokens"`
}

func AvailableTokens(token Token) PayTokenListOutput {
	var arrTokens []Token = []Token{
		Token{"ETH", "ETH"},
		Token{"BTC", "BTC"},
		Token{"USDT", "USDT"},
	}
	arrTokens = append(arrTokens, token)
	var payTokens PayTokenListOutput = PayTokenListOutput{arrTokens}
	return payTokens
}
