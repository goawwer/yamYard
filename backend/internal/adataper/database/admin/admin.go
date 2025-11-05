package admin

import (
	"context"

	"github.com/goawwer/yamyard/internal/adataper/database"
	"github.com/goawwer/yamyard/internal/domain"
	"github.com/google/uuid"
)

type AdminRepository struct {
	r *database.Db
}

func NewAdminRepository(db *database.Db) *AdminRepository {
	return &AdminRepository{r: db}
}

func (a *AdminRepository) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	query := `
        SELECT id, email, username, image_url, bio, profile_status, is_admin, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
    `
	rows, err := a.r.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u domain.User
		err := rows.Scan(
			&u.ID, &u.Email, &u.Username, &u.ImageURL,
			&u.Bio, &u.ProfileStatus, &u.IsAdmin,
			&u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}

func (a *AdminRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := a.r.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (a *AdminRepository) GetAllRecipes(ctx context.Context) ([]*domain.Recipe, error) {
	var recipes []*domain.Recipe
	query := `
        SELECT r.id, r.author_id, u.username, r.title, r.description, 
               r.cooking_time, r.difficulty, r.image_url, r.created_at, r.likes_count
        FROM recipes r
        JOIN users u ON r.author_id = u.id
        ORDER BY r.created_at DESC
    `
	rows, err := a.r.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rec domain.Recipe
		var username string
		err := rows.Scan(
			&rec.ID, &rec.AuthorID, &username,
			&rec.Title, &rec.Description, &rec.CookingTime,
			&rec.Difficulty, &rec.ImageURL, &rec.CreatedAt, &rec.LikesCount,
		)
		if err != nil {
			return nil, err
		}
		rec.AuthorUsername = username
		recipes = append(recipes, &rec)
	}
	return recipes, rows.Err()
}

func (a *AdminRepository) DeleteRecipe(ctx context.Context, id uuid.UUID) error {
	_, err := a.r.ExecContext(ctx, `DELETE FROM recipes WHERE id = $1`, id)
	return err
}
