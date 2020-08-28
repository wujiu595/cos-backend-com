package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
	"cos-backend-com/src/libs/sdk/eth"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type BountyStatus int

const (
	BountyStatusOpen       BountyStatus = iota
	BountyStatusInProgress BountyStatus = iota
	BountyStatusClosed     BountyStatus = iota
)

type BountyType string

const (
	BountyTypeContest     BountyType = "contest"
	BountyTypeCooperative BountyType = "cooperative"
)

func (b BountyType) Validate() bool {
	switch b {
	case BountyTypeContest, BountyTypeCooperative:
		return true
	}
	return false
}

type BountyHunterRelStatus int

const (
	BountyStatusStartWork BountyHunterRelStatus = iota
	BountyStatusSubmitted BountyHunterRelStatus = iota
	BountyStatusPaid      BountyHunterRelStatus = iota
	BountyStatusQuited    BountyHunterRelStatus = iota
)

type CreateBountyInput struct {
	Id              flake.ID       `json:"id" validate:"required"`
	Title           string         `json:"title" validate:"required"`
	TxId            string         `json:"txId" validate:"required"`
	Type            string         `json:"type" validate:"func=self.Validate"`
	Keywords        []string       `json:"keywords"`
	ContactEmail    string         `json:"contactEmail" validate:"required"`
	Intro           string         `json:"intro" validate:"required"`
	DescriptionAddr string         `json:"descriptionAddr"`
	Duration        int32          `json:"duration" validate:"required"`
	Payments        types.JSONText `json:"payments" validate:"required"`
}

type BountyOutput struct {
	Id      flake.ID `json:"id" db:"id"`
	Startup struct {
		Id   flake.ID `json:"id" db:"id"`
		Name string   `json:"name" db:"name"`
	} `json:"startupId" db:"startup_id"`
	Uid                 flake.ID       `json:"uid" db:"uid"`
	ContactEmail        string         `json:"contactEmail" db:"contact_email"`
	BlockAddr           string         `json:"blockAddr" db:"block_addr"`
	DescriptionAddr     string         `json:"descriptionAddr" db:"description_addr"`
	DescriptionFileAddr string         `json:"descriptionFileAddr" db:"description_file_addr"`
	Payments            types.JSONText `json:"payments" validate:"required"`
	Hunters             []struct {
		HunterId    flake.ID              `json:"hunterId" db:"hunter_id"`       // hunter_id
		Name        string                `json:"name" db:"name"`                // name
		Status      BountyHunterRelStatus `json:"status" db:"status"`            // status
		StartedAt   *time.Time            `json:"startedAt" db:"started_at"`     // started_at
		SubmittedAt *time.Time            `json:"submittedAt" db:"submitted_at"` // submitted_at
		QuitedAt    *time.Time            `json:"quitedAt" db:"quited_at"`       // quited_at
		PaidAt      *time.Time            `json:"paidAt" db:"paid_at"`           // paid_at
		PaidTokens  types.JSONText        `json:"paidTokens" db:"paid_tokens"`   // paid_tokens
	} `json:"hunters" db:"hunters"`
	Status           BountyStatus         `json:"status" db:"status"`
	TransactionState eth.TransactionState `json:"transactionState" db:"transaction_state"`
}

type ListBountiesInput struct {
	Keyword string `param:"keyword"`
	pagination.ListRequest
}

type ListBountiesResult struct {
	pagination.ListResult
	Result []BountyOutput `json:"result"`
}
