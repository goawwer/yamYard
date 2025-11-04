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

func (p *ProfileRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*domain.User, error) {
	var user domain.User

	query := `
    SELECT
        id,
        email,
        username,
        hashed_password,
        bio,
        image_url,
        profile_status,
        is_admin,
        created_at,
        updated_at
    FROM users
    WHERE id = $1
`

	err := p.r.QueryRowContext(ctx, query, userId).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.HashedPassword,
		&user.Bio,
		&user.ImageURL,
		&user.ProfileStatus,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return &user, err
}

func (p *ProfileRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	err := p.r.QueryRowContext(ctx, query, id).Scan(&exists)

	return exists, err
}

func (p *ProfileRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
        UPDATE users
        SET 
            username = $1,
            profile_status = $2,
            bio = $3,
			image_url = $4,
            updated_at = current_timestamp
        WHERE id = $5
        RETURNING 
            id, username, email, hashed_password, 
            bio, image_url, profile_status, is_admin, 
            created_at, updated_at
    `

	var updated domain.User

	err := p.r.QueryRowContext(ctx, query,
		user.Username,
		user.ProfileStatus,
		user.Bio,
		user.ImageURL,
		user.ID,
	).Scan(
		&updated.ID,
		&updated.Username,
		&updated.Email,
		&updated.HashedPassword,
		&updated.Bio,
		&updated.ImageURL,
		&updated.ProfileStatus,
		&updated.IsAdmin,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updated, nil
}
