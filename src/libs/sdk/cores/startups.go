package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
	"cos-backend-com/src/libs/sdk/eth"
)

type StartUpState int

const (
	StartUpStateCreating      StartUpState = 0
	StartUpStateCreated       StartUpState = 1
	StartUpStateConfirmFailed StartUpState = 2
	StartUpStateFailed        StartUpState = 3
	StartUpStateHasSetting    StartUpState = 4
)

type CreateStartupInput struct {
	CreateStartupRevisionInput
}

type UpdateStartupInput struct {
	CreateStartupRevisionInput
}

type CreateStartupRevisionInput struct {
	Name            string   `json:"name"`
	Mission         *string  `json:"mission"`
	Logo            string   `json:"logo"`
	DescriptionAddr string   `json:"descriptionAddr"`
	CategoryId      flake.ID `json:"categoryId"`
	TxId            string   `json:"txId"`
}
type StartupIdResult struct {
	Id flake.ID `json:"id" db:"id"`
}

type StartUpResult struct {
	Id              flake.ID             `json:"id" db:"id"`
	Name            string               `json:"name" db:"name"`
	Mission         *string              `json:"mission" db:"mission"`
	Logo            string               `json:"logo" db:"logo"`
	DescriptionAddr string               `json:"descriptionAddr" db:"description_addr"`
	Category        CategoriesResult     `json:"category" db:"category"`
	State           eth.TransactionState `json:"state" db:"state"`
	SettingState    eth.TransactionState `json:"settingState" db:"setting_state"`
	IsIRO           bool                 `json:"isIRO" db:"is_iro"`
}

type ListStartupsInput struct {
	CategoryId flake.ID `param:"categoryId"`
	IsIRO      bool     `param:"isIRo"`
	Keyword    string   `param:"keyword"`
	pagination.ListRequest
}

type ListStartupsResult struct {
	pagination.ListResult
	Result []StartUpResult `json:"result"`
}
