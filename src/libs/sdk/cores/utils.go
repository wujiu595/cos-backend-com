package cores

import (
	"cos-backend-com/src/common/flake"
)

type PrepareIdOutput struct {
	Id flake.ID `json:"id"`
}

type Token struct {
	name   string
	symbol string
}

func AvailableTokens(name string, symbol string) []Token {
	var arrTokens []Token = []Token{
		Token{"ETH", "ETH"},
		Token{"BTC", "BTC"},
		Token{"USDT", "USDT"},
	}
	arrTokens = append(arrTokens, Token{name, symbol})

	return arrTokens
}
