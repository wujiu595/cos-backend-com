package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
)

type CategorySource string

const CategorySourceStartUp = "startup"

type CategoriesInput struct {
	Name   string         `json:"name"`
	Code   string         `json:"code"`
	Source CategorySource `json:"source"`
}

type CategoriesResult struct {
	Id     flake.ID       `json:"id" db:"id"`
	Name   string         `json:"name" db:"name"`
	Code   string         `json:"code" db:"code"`
	Source CategorySource `json:"source" db:"source"`
}

type ListCategoriesInput struct {
	pagination.ListRequest
	ShowDeleted bool           `param:"showDeleted"`
	Source      CategorySource `param:"source"`
}

type ListCategoriesResult struct {
	pagination.ListResult
	Result []CategoriesResult `json:"result"`
}
