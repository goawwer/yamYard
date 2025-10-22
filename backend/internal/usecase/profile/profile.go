package profile

import (
	"context"
	"fmt"
	"time"

	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/internal/dto"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (p *ProfileService) SignUp(ctx context.Context, input dto.SignUpUserInput) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := domain.User{
		Email:          input.Email,
		Username:       input.Username,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now(),
	}

	if err := user.Validate(); err != nil {
		return fmt.Errorf("failed to validate user input arguments: %w", err)
	}

	if err := p.repo.Create(ctx, &user); err != nil {
		return fmt.Errorf("failed to add user to database: %w", err)
	}

	return nil
}

func (p *ProfileService) Login(ctx context.Context, input dto.LoginInput) (uuid.UUID, error) {
	output, err := p.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return uuid.Nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(output.HashedPassword), []byte(input.Password))
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid credentials: %w", err)
	}

	return output.ID, nil
}

func (p *ProfileService) Exists(ctx context.Context, userId uuid.UUID) (bool, error) {
	output, err := p.repo.ExistsByID(ctx, userId)
	if err != nil {
		return false, fmt.Errorf("user not found")
	}

	return output, nil
}
