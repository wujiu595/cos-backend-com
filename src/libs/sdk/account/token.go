package account

import (
	"cos-backend-com/src/common/flake"
	"time"
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

type AccessTokensResult struct {
	Id        flake.ID  `json:"id" db:"id"`                // id
	UId       flake.ID  `json:"uid" db:"uid"`              // uid
	Token     string    `json:"token" db:"token"`          // token
	Refresh   string    `json:"refresh" db:"refresh"`      // refresh
	Key       string    `json:"key" db:"key"`              // key
	Secret    string    `json:"secret" db:"secret"`        // secret
	CreatedAt time.Time `json:"createdAt" db:"created_at"` // created_at
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"` // updated_at
}
