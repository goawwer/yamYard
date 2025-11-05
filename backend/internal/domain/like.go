package domain

import (
	"time"

	"github.com/google/uuid"
)

// internal/domain/like.go
type Like struct {
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	RecipeID  uuid.UUID `db:"recipe_id" json:"recipe_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
