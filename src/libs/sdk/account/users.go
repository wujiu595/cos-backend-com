package account

import (
	"cos-backend-com/src/common/flake"
	"time"
)

const DefaultNoncePrefix = "The nonce for login comunion is:"

type Signature string

type UsersModel struct {
	Id            flake.ID  `json:"id" db:"id"`                        // id (PK)
	PublicKey     string    `json:"PublicKey" db:"public_key"`         // public_key
	Nonce         string    `json:"nonce" db:"nonce"`                  // nonce
	PublicSecret  string    `json:"publicSecret" db:"public_secret"`   // public_secret
	PrivateSecret string    `json:"privateSecret" db:"private_secret"` // private_secret
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`         // created_at
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`         // updated_at
	IsHunter      bool      `json:"isHunter" db:"is_hunter"`           // is_hunter
}

type GetNonceInput struct {
	PublicAddr string `json:"publicKey"`
}

type GetNonceOutput struct {
	Nonce string `json:"nonce"`
}

type LoginInput struct {
	PublicKey string `json:"publicKey" validate:"func "`
	Signature string `json:"signature"`
}

type UserResult struct {
	Id        flake.ID `json:"id" db:"id"`                // id (PK)
	PublicKey string   `json:"publicKey" db:"public_key"` // wallet_addr
	IsHunter  bool     `json:"isHunter" db:"is_hunter"`   // is_hunter
}
