package like

import (
	"context"

	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/google/uuid"
)

type LikeRepo interface {
	Toggle(ctx context.Context, userID, recipeID uuid.UUID) (bool, int, error)
	LikedRecipes(ctx context.Context, userID uuid.UUID, input helpers.FilterAndSortingParameters) ([]*domain.Recipe, error)
}

type recipeRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Recipe, error)
}

// internal/usecase/like/service.go
type LikeService struct {
	likeRepo   LikeRepo
	recipeRepo recipeRepo
}

func NewLikeService(l LikeRepo, r recipeRepo) *LikeService {
	return &LikeService{likeRepo: l, recipeRepo: r}
}
