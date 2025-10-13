package user

import (
	"context"
	"fmt"
	"time"

	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserService) SignUp(ctx context.Context, input dto.SignUpUserInput) (dto.SinUpUserOutput, error) {
	var output dto.SignUpUserOutput

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return output, fmt.Errorf("failed to hash password")
	}

	user := domain.User {
		Email: input.Email,
		Username: input.Username,
		HashedPassword: hashedPassword,
		CreatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return output, fmt.Errorf("failed to validate user input arguments")
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return output, err
	}

	return dto.SignUpUserOutput{
		ID: user.ID,
	}, nil
}

