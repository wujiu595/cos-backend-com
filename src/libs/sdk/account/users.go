package account

import (
	"cos-backend-com/src/common/flake"
)

type UsersResult struct {
	Id         flake.ID `json:"id" db:"id"`                  // id (PK)
	WalletAddr string   `json:"walletAddr" db:"wallet_addr"` // wallet_addr
	IsHunter   bool     `json:"isHunter" db:"is_hunter"`     // is_hunter
}
