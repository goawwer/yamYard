package profile

import (
	"context"

	"github.com/goawwer/yamyard/internal/adataper/database"
	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/internal/dto"
)

type ProfileRepository struct {
	r *database.Db
}

func NewProfileRepository(db *database.Db) *ProfileRepository {
	return &ProfileRepository{
		r: db,
	}
}

func (p *ProfileRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, username, hashed_password, created_at)
		VALUES ($1, $2, $3, $4)
	`

	return p.r.QueryRowContext(
		ctx, query,
		user.Email, user.Username, user.HashedPassword, user.CreatedAt,
	).Err()
}

func (p *ProfileRepository) GetUserByEmail(ctx context.Context, email string) (dto.LoginOutput, error) {
	var output dto.LoginOutput

	query := `
		SELECT * FROM users
		WHERE email = $1
	`
	err := p.r.QueryRowContext(ctx, query, email).Scan(&output.ID)

	return output, err
}
