package account

import "time"

type UpdateHunterInput struct {
	Name            string   `json:"name"`
	Skills          []string `json:"skills"`
	About           string   `json:"about"`
	DescriptionAddr string   `json:"descriptionAddr"`
	Email           string   `json:"email"`
}

type HunterResult struct {
	Name            string    `json:"name" db:"name"`                        // name
	Skills          []string  `json:"skills" db:"skills"`                    // skills
	About           string    `json:"about" db:"about"`                      // about
	DescriptionAddr string    `json:"descriptionAddr" db:"description_addr"` // description_addr
	Email           string    `json:"email" db:"email"`                      // email
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`             //created_at
}
