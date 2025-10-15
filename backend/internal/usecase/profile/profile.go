package profile

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/internal/dto"
	"github.com/goawwer/yamyard/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

func (p *ProfileService) SignUp(ctx context.Context, input dto.SignUpUserInput) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	user := domain.User{
		Email:          input.Email,
		Username:       input.Username,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now(),
	}

	if err := user.Validate(); err != nil {
		return fmt.Errorf("failed to validate user input arguments")
	}

	if err := p.repo.Create(ctx, &user); err != nil {
		return err
	}

	return nil
}

func (p *ProfileService) Login(ctx context.Context, input dto.LoginInput) (string, error) {
	output, err := p.repo.GetUserByEmail(ctx, input.Email)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("user not found")
	} else if err != nil {
		return "", fmt.Errorf("database internal error")
	}

	token, err := middleware.GenerateToken(output.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}

	return token, nil
}
