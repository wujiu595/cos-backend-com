package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
)

type StartUpState int

const (
	StartUpStateCreating      StartUpState = 0
	StartUpStateCreated       StartUpState = 1
	StartUpStateConfirmFailed StartUpState = 2
	StartUpStateFailed        StartUpState = 3
	StartUpStateHasSetting    StartUpState = 4
)

type CreateStartUpsInput struct {
	Name            string       `json:"name"`
	Mission         *string      `json:"mission"`
	Logo            string       `json:"logo"`
	TxId            string       `json:"txId"`
	BlockNum        *flake.ID    `json:"blockNum"`
	DescriptionAddr string       `json:"descriptionAddr"`
	CategoryId      flake.ID     `json:"categoryId"`
	State           StartUpState `json:"state"`
	IsIRO           bool         `json:"isIRO"`
}

type StartUpsResult struct {
	Id              flake.ID         `json:"id" db:"id"`
	Name            string           `json:"name" db:"name"`
	Mission         *string          `json:"mission" db:"mission"`
	Logo            string           `json:"logo" db:"logo"`
	TxId            string           `json:"txId" db:"tx_id"`
	BlockNum        *int64           `json:"blockNum" db:"block_num"`
	DescriptionAddr string           `json:"descriptionAddr" db:"description_addr"`
	Category        CategoriesResult `json:"category" db:"category"`
	State           StartUpState     `json:"state" db:"state"`
	IsIRO           bool             `json:"isIRO" db:"is_iro"`
}

type ListStartUpsInput struct {
	pagination.ListRequest
}

type ListStartUpsResult struct {
	pagination.ListResult
	Result []StartUpsResult `json:"result"`
}
