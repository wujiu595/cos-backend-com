package cores

import (
	"cos-backend-com/src/common/flake"
)

type StartupSettingResult struct {
	Id flake.ID `json:"id" db:"id"` // id
}

type UpdateStartupSettingInput struct {
	TxId        string `json:"txId"`
	TokenName   string `json:"tokenName"`
	TokenSymbol string `json:"tokenSymbol"`
	TokenAddr   string `json:"tokenAddr"`
	WalletAddrs []struct {
		Name string `json:"name"`
		Addr string `json:"addr"`
	} `json:"walletAddrs"`
	Type                   string   `json:"type"`
	VoteTokenLimit         int64    `json:"voteTokenLimit"`
	VoteAssignAddrs        []string `json:"voteAssignAddrs"`
	VoteSupportPercent     int      `json:"voteSupportPercent"`
	VoteMinApprovalPercent int      `json:"voteMinApprovalPercent"`
	VoteMinDurationHours   int64    `json:"voteMinDurationHours"`
	VoteMaxDurationHours   int64    `json:"voteMaxDurationHours"`
}
