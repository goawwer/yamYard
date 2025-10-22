package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `db:"id" json:"id"`
	Email          string    `db:"email" json:"email"`
	Username       string    `db:"username" json:"username"`
	HashedPassword string    `db:"hashed_password" json:"password"`
	Bio            string    `db:"bio" json:"bio"`
	ProfileStatus  string    `db:"profile_status" json:"profile_status"`
	IsAdmin        bool      `db:"is_admin" json:"is_admin"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
type UserEntity struct {
	ID             uuid.UUID `db:"id"`
	Email          string    `db:"email"`
	Username       string    `db:"username"`
	HashedPassword string    `db:"hashed_password"`
}

func (u User) Validate() error {
	if u.Email == "" || u.Username == "" {
		return ErrCredentialsRequired
	}

	if len(u.Username) < 5 || len(u.Username) > 64 {
		return ErrInvalidUsername
	}

	return nil
}
