package profile

import (
	"context"

	"github.com/goawwer/yamyard/internal/domain"
	"github.com/google/uuid"
)

type ProfileRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.UserEntity, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
}

type ProfileService struct {
	repo ProfileRepository
}

func NewProfileService(repo ProfileRepository) *ProfileService {
	return &ProfileService{
		repo: repo,
	}
}
