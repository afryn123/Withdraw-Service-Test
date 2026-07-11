package application_test

import (
	"context"
	"testing"

	"github.com/afryn123/withdraw-service-test/internal/mocks"
	userdomain "github.com/afryn123/withdraw-service-test/internal/user/domain"
	"github.com/afryn123/withdraw-service-test/internal/wallet/application"
	walletdomain "github.com/afryn123/withdraw-service-test/internal/wallet/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestServiceFindBalanceByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository := mocks.NewMockWalletRepository(ctrl)
	userID, walletID := uuid.New(), uuid.New()
	repository.EXPECT().FindByUserID(gomock.Any(), userID).Return(&walletdomain.Wallet{
		ID: walletID, UserID: userID, Balance: 125_000,
		User: userdomain.User{ID: userID, Name: "Afriyan"},
	}, nil)

	result, err := application.NewService(repository).FindBalanceByUserID(context.Background(), userID)
	require.NoError(t, err)
	require.Equal(t, walletID, result.WalletID)
	require.EqualValues(t, 125_000, result.Balance)
	require.Equal(t, "Afriyan", result.User.Name)
}
