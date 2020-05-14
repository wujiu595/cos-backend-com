package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
)

type CategorySource string

const CategorySourceStartUp = "startup"

func (c CategorySource) Validate() bool {
	switch c {
	case CategorySourceStartUp, "":
		return true
	}
	return false
}

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
	Source      CategorySource `param:"source" validate:"func=self.Validate"`
}

type ListCategoriesResult struct {
	pagination.ListResult
	Result []CategoriesResult `json:"result"`
}
