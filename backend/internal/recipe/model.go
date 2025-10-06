package recipe

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"recipe_name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
