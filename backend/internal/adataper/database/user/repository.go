package user

import (
	"context"

	"github.com/goawwer/yamyard/internal/adataper/database"
	"github.com/goawwer/yamyard/internal/domain"
)


type UserRepository struct {
	r *database.Db
}	

func NewUserRepository(db *database.Db) *UserRepository {
	return &UserRepository{
		r: db,	
	}
}


func (u *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, username, hashed_password, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	return u.r.QueryRowContext(
		ctx, query,
		user.Email, user.Username, user.HashedPassword, user.CreatedAt,
		).Scan(&user.ID)
} 
