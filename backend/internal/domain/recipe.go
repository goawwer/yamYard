package domain

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID              uuid.UUID `db:"id" json:"id"`
	AuthorUsername  string    `db:"-" json:"author_username"`
	AuthorID        uuid.UUID `db:"author_id" json:"author_id"`
	AuthorAvatarURL *string   `db:"-" json:"author_avatar_url"`
	Title           string    `db:"title" json:"title"`
	Description     string    `db:"description" json:"description"`
	Ingredients     string    `db:"ingredients" json:"ingredients"`
	CookingTime     *int      `db:"cooking_time" json:"cooking_time"`
	Difficulty      *string   `db:"difficulty" json:"difficulty"`
	ImageURL        *string   `db:"image_url" json:"image_url"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
	LikesCount      int       `json:"likes_count"`
	IsLiked         bool      `json:"is_liked"`
}
