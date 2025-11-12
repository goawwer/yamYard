package like

import (
	"context"
	"fmt"
	"strings"

	"github.com/goawwer/yamyard/internal/adataper/database"
	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/google/uuid"
)

// internal/adataper/repository/like.go
type LikeRepository struct{ r *database.Db }

func NewLikeRepository(db *database.Db) *LikeRepository { return &LikeRepository{r: db} }

// Toggle = INSERT or DELETE
func (l *LikeRepository) Toggle(ctx context.Context, userID, recipeID uuid.UUID) (bool, int, error) {
	var exists bool
	err := l.r.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND recipe_id = $2)`, userID, recipeID).Scan(&exists)
	if err != nil {
		return false, 0, err
	}

	if exists {
		_, err = l.r.ExecContext(ctx, `DELETE FROM likes WHERE user_id = $1 AND recipe_id = $2`, userID, recipeID)
	} else {
		_, err = l.r.ExecContext(ctx, `INSERT INTO likes (user_id, recipe_id) VALUES ($1, $2)`, userID, recipeID)
	}
	if err != nil {
		return false, 0, err
	}

	var count int
	err = l.r.QueryRowContext(ctx, `SELECT likes_count FROM recipes WHERE id = $1`, recipeID).Scan(&count)
	if err != nil {
		return false, 0, err
	}

	return !exists, count, nil
}

// IsLiked – for the heart icon
func (l *LikeRepository) IsLiked(ctx context.Context, userID, recipeID uuid.UUID) (bool, error) {
	var yes bool
	err := l.r.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM likes WHERE user_id=$1 AND recipe_id=$2)`,
		userID, recipeID).Scan(&yes)
	return yes, err
}

// LikedRecipes – for the “Liked” tab
func (l *LikeRepository) LikedRecipes(ctx context.Context, userID uuid.UUID, input helpers.FilterAndSortingParameters) ([]*domain.Recipe, error) {
	var recipes []*domain.Recipe

	query := `
        SELECT 
            r.id, r.author_id, u.username, u.image_url as author_avatar_url,
            r.title, r.description, r.ingredients, r.cooking_time, 
            r.difficulty, r.image_url, r.created_at, r.updated_at,
            r.likes_count,
            EXISTS(SELECT 1 FROM likes ll WHERE ll.recipe_id = r.id AND ll.user_id = $1) AS is_liked
        FROM likes l
        JOIN recipes r ON l.recipe_id = r.id
        JOIN users u ON r.author_id = u.id
        WHERE l.user_id = $2
    `

	args := []any{userID, userID} // $1 for subquery, $2 for WHERE
	where := []string{}

	if input.Key != "" && input.Value != "" {
		where = append(where, fmt.Sprintf("r.%s ILIKE $%d", input.Key, len(args)+1))
		args = append(args, "%"+input.Value+"%")
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	// Sorting (same as GetAll)
	sortField := "r.created_at"
	if input.Sort != "" {
		sortField = "r." + input.Sort
	}
	order := "DESC"
	if strings.ToLower(input.Order) == "asc" {
		order = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortField, order)

	rows, err := l.r.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query liked recipes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rec domain.Recipe
		var username string
		var authorAvatarURL *string
		var isLiked bool

		err := rows.Scan(
			&rec.ID, &rec.AuthorID, &username, &authorAvatarURL,
			&rec.Title, &rec.Description, &rec.Ingredients, &rec.CookingTime,
			&rec.Difficulty, &rec.ImageURL, &rec.CreatedAt, &rec.UpdatedAt,
			&rec.LikesCount, &isLiked,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan liked recipe: %w", err)
		}

		rec.AuthorUsername = username
		rec.AuthorAvatarURL = authorAvatarURL
		rec.IsLiked = isLiked
		recipes = append(recipes, &rec)
	}

	return recipes, rows.Err()
}
