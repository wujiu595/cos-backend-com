package account

import (
	"cos-backend-com/src/common/flake"

	"github.com/dgrijalva/jwt-go"
)

const (
	TokenTypeRefreshToken TokenType = "rt"
	TokenTypeAccessToken  TokenType = "at"
)

type TokenType string

func (p TokenType) Valid() bool {
	switch p {
	case TokenTypeRefreshToken, TokenTypeAccessToken:
	default:
		return false
	}
	return true
}

func (p TokenType) String() string {
	return string(p)
}

type TokenInfo struct {
	Uid flake.ID `json:"uid"`
}

type RefreshTokenClaims struct {
	TokenClaims
}

type AccessTokenClaims struct {
	TokenClaims
}

type TokenClaims struct {
	TokenInfo
	jwt.StandardClaims
}
