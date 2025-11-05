package like

import (
	"context"

	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/google/uuid"
)

func (s *LikeService) Toggle(ctx context.Context, userID, recipeID uuid.UUID) (bool, int, error) {
	return s.likeRepo.Toggle(ctx, userID, recipeID)
}

func (s *LikeService) LikedRecipes(ctx context.Context, userID uuid.UUID, input helpers.FilterAndSortingParameters) ([]*domain.Recipe, error) {
	return s.likeRepo.LikedRecipes(ctx, userID, input)
}
