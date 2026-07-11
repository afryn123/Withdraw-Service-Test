package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/afryn123/withdraw-service-test/internal/auth/application"
	authdomain "github.com/afryn123/withdraw-service-test/internal/auth/domain"
	"github.com/afryn123/withdraw-service-test/internal/mocks"
	userdomain "github.com/afryn123/withdraw-service-test/internal/user/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestServiceLogin(t *testing.T) {
	t.Run("returns token for valid credentials", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		users := mocks.NewMockUserRepository(ctrl)
		passwords := mocks.NewMockPasswordVerifier(ctrl)
		tokens := mocks.NewMockTokenManager(ctrl)
		userID := uuid.New()
		users.EXPECT().FindByEmail(gomock.Any(), "a@example.com").Return(&userdomain.User{ID: userID, Password: "hash"}, nil)
		passwords.EXPECT().Verify("secret12", "hash").Return(true)
		tokens.EXPECT().Generate(userID).Return("jwt", nil)

		token, err := application.NewService(users, passwords, tokens).Login(context.Background(), "a@example.com", "secret12")
		require.NoError(t, err)
		require.Equal(t, "jwt", token)
	})

	t.Run("hides repository details behind invalid credentials", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		users := mocks.NewMockUserRepository(ctrl)
		users.EXPECT().FindByEmail(gomock.Any(), "none@example.com").Return(nil, errors.New("not found"))

		_, err := application.NewService(users, mocks.NewMockPasswordVerifier(ctrl), mocks.NewMockTokenManager(ctrl)).
			Login(context.Background(), "none@example.com", "secret12")
		require.ErrorIs(t, err, authdomain.ErrInvalidCredentials)
	})
}
