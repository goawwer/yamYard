package profile

import (
	"context"

	"github.com/goawwer/yamyard/internal/adataper/database"
	"github.com/goawwer/yamyard/internal/domain"
	"github.com/google/uuid"
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

func (p *ProfileRepository) GetUserByEmail(ctx context.Context, email string) (domain.UserEntity, error) {
	var output domain.UserEntity

	query := `
		SELECT id, email, username, hashed_password FROM users
		WHERE email = $1
	`
	err := p.r.QueryRowContext(ctx, query, email).Scan(&output.ID, &output.Email, &output.Username, &output.HashedPassword)

	return output, err
}

func (p *ProfileRepository) GetUserById(ctx context.Context, userId int) error {
	query := `
		SELECT * FROM users
		WHERE id = $1
	`
	return p.r.QueryRowContext(ctx, query, userId).Scan(&userId)
}

func (p *ProfileRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	err := p.r.QueryRowContext(ctx, query, id).Scan(&exists)

	return exists, err
}
