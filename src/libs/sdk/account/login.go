package account

type LoginInput struct {
	WalletAddr string `json:"walletAddr"`
}

type LoginUserResult struct {
	UsersResult
	PublicSecret  string `json:"publicSecret" db:"public_secret"`
	PrivateSecret string `json:"privateSecret" db:"private_secret"`
}
