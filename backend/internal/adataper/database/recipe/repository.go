package recipe

import (
	"context"
	"fmt"
	"strings"

	"github.com/goawwer/yamyard/internal/adataper/database"
	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/google/uuid"
)

type RecipeRepository struct {
	r *database.Db
}

func NewRecipeRepository(db *database.Db) *RecipeRepository {
	return &RecipeRepository{
		r: db,
	}
}

func (r *RecipeRepository) Create(ctx context.Context, recipe *domain.Recipe) error {
	query := `
		INSERT INTO recipes (
			author_id, title, description, ingredients, cooking_time, difficulty, image_url
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	return r.r.QueryRowContext(
		ctx,
		query,
		recipe.AuthorID,
		recipe.Title,
		recipe.Description,
		recipe.Ingredients,
		recipe.CookingTime,
		recipe.Difficulty,
		recipe.ImageURL,
	).Scan(&recipe.ID)
}

func (r *RecipeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Recipe, error) {
	var recipe domain.Recipe

	query := `
		SELECT * FROM recipes 
		WHERE id = $1
	`

	err := r.r.QueryRowContext(ctx, query, id).Scan(
		&recipe.ID,
		&recipe.AuthorID,
		&recipe.Title,
		&recipe.Description,
		&recipe.Ingredients,
		&recipe.CookingTime,
		&recipe.Difficulty,
		&recipe.ImageURL,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
	)

	return &recipe, err
}

func (r *RecipeRepository) GetAll(ctx context.Context, input helpers.FilterAndSortingParameters) ([]*domain.Recipe, error) {
	var recipes []*domain.Recipe

	// Base query
	query := `
  		SELECT 
   		 	r.id, r.author_id, u.username, u.image_url as author_avatar_url,
    		r.title, r.description, r.ingredients, r.cooking_time, 
    		r.difficulty, r.image_url, r.created_at, r.updated_at
  		FROM recipes r
  		JOIN users u ON r.author_id = u.id
	`

	// arguments for WHERE placeholders
	args := []any{}
	where := []string{}

	if input.Username != "" {
		where = append(where, fmt.Sprintf("u.username ILIKE $%d", len(args)+1))
		args = append(args, "%"+input.Username+"%")
	}

	if input.Key != "" && input.Value != "" {
		where = append(where, fmt.Sprintf("r.%s ILIKE $%d", input.Key, len(args)+1))
		args = append(args, "%"+input.Value+"%")
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	sortField := "r.created_at"
	if input.Sort != "" {
		sortField = "r." + input.Sort
	}

	order := "DESC"
	if strings.ToLower(input.Order) == "asc" {
		order = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortField, order)
	rows, err := r.r.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rec domain.Recipe
		var username string
		var authorAvatarURL *string
		err := rows.Scan(
			&rec.ID,
			&rec.AuthorID,
			&username,
			&authorAvatarURL,
			&rec.Title,
			&rec.Description,
			&rec.Ingredients,
			&rec.CookingTime,
			&rec.Difficulty,
			&rec.ImageURL,
			&rec.CreatedAt,
			&rec.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %w", err)
		}

		rec.AuthorUsername = username
		recipes = append(recipes, &rec)
	}

	return recipes, rows.Err()
}

func (r *RecipeRepository) Update(ctx context.Context, recipe *domain.Recipe) error {
	query := `
		UPDATE recipes
		SET title = $1,
		    description = $2,
		    ingredients = $3,
		    cooking_time = $4,
		    difficulty = $5,
		    image_url = $6,
		    updated_at = current_timestamp
		WHERE id = $7
	`

	_, err := r.r.ExecContext(
		ctx,
		query,
		recipe.Title,
		recipe.Description,
		recipe.Ingredients,
		recipe.CookingTime,
		recipe.Difficulty,
		recipe.ImageURL,
		recipe.ID,
	)
	return err
}

func (r *RecipeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM recipes WHERE id = $1`
	_, err := r.r.ExecContext(ctx, query, id)
	return err
}
