package profile

import (
	"context"

	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/google/uuid"
)

type ProfileRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.UserEntity, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (*domain.User, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	Update(ctx context.Context, newUser *domain.User) (*domain.User, error)
	GetAllUsers(ctx context.Context, input helpers.FilterAndSortingParameters) ([]*domain.User, error)
}

type ProfileService struct {
	repo ProfileRepository
}

func NewProfileService(repo ProfileRepository) *ProfileService {
	return &ProfileService{
		repo: repo,
	}
}
