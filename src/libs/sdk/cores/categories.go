package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
)

type CategoriesInput struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Source string `json:"source"`
}

type CategoriesResult struct {
	Id     flake.ID `json:"id" db:"id"`
	Name   string   `json:"name" db:"name"`
	Code   string   `json:"code" db:"code"`
	Source string   `json:"source" db:"source"`
}

type ListCategoriesInput struct {
	pagination.ListRequest
	ShowDeleted bool `param:"showDeleted"`
}

type ListCategoriesResult struct {
	pagination.ListResult
	Result []CategoriesResult `json:"result"`
}
