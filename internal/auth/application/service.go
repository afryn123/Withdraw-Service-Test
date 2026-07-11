package application

import (
	"context"

	authdomain "github.com/afryn123/withdraw-service-test/internal/auth/domain"
	userdomain "github.com/afryn123/withdraw-service-test/internal/user/domain"
)

type Service struct {
	users    userdomain.Repository
	password authdomain.PasswordVerifier
	tokens   authdomain.TokenManager
}

func NewService(users userdomain.Repository, password authdomain.PasswordVerifier, tokens authdomain.TokenManager) *Service {
	return &Service{users: users, password: password, tokens: tokens}
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil || !s.password.Verify(password, user.Password) {
		return "", authdomain.ErrInvalidCredentials
	}
	return s.tokens.Generate(user.ID)
}
