package middleware

import (
	"context"

	"github.com/goawwer/yamyard/internal/adataper/database/profile"
	"github.com/google/uuid"
)

type AuthService struct {
	repo profile.ProfileRepository
}

func NewAuthService(repo profile.ProfileRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repo.ExistsByID(ctx, id)
}
