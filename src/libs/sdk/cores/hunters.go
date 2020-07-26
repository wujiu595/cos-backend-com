package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
)

type CreateHunterInput struct {
	Name            string   `json:"name"`
	Skills          []string `json:"skills"`
	About           string   `json:"about"`
	DescriptionAddr string   `json:"descriptionAddr"`
	Email           string   `json:"email"`
}

type UpdateHunterInput struct {
	Name            string   `json:"name"`
	Skills          []string `json:"skills"`
	About           string   `json:"about"`
	DescriptionAddr string   `json:"descriptionAddr"`
	Email           string   `json:"email"`
}

type HunterResult struct {
	Id              flake.ID `json:"id" db:"id"`                            // id (PK)
	UserId          flake.ID `json:"userId" db:"user_id"`                   // user_id
	Name            string   `json:"name" db:"name"`                        // name
	Skills          []string `json:"skills" db:"skills"`                    // skills
	About           string   `json:"about" db:"about"`                      // about
	DescriptionAddr string   `json:"descriptionAddr" db:"description_addr"` // description_addr
	Email           string   `json:"email" db:"email"`                      // email
}

type ListHuntersInput struct {
	pagination.ListRequest
}

type ListHuntersResult struct {
	pagination.ListResult
	Result []HunterResult `json:"result"`
}
