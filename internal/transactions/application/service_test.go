package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/afryn123/withdraw-service-test/internal/mocks"
	"github.com/afryn123/withdraw-service-test/internal/transactions/application"
	transactiondomain "github.com/afryn123/withdraw-service-test/internal/transactions/domain"
	userdomain "github.com/afryn123/withdraw-service-test/internal/user/domain"
	walletdomain "github.com/afryn123/withdraw-service-test/internal/wallet/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func executeTransaction(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) }

func TestServiceWithdraw(t *testing.T) {
	t.Run("locks wallet, updates balance, and records transaction", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		wallets := mocks.NewMockWalletRepository(ctrl)
		transactions := mocks.NewMockTransactionRepository(ctrl)
		tx := mocks.NewMockTransactionManager(ctrl)
		codes := mocks.NewMockCodeGenerator(ctrl)
		userID, walletID, transactionID := uuid.New(), uuid.New(), uuid.New()
		remark := "ATM"
		wallet := &walletdomain.Wallet{
			ID: walletID, UserID: userID, Balance: 100_000,
			User: userdomain.User{ID: userID, Name: "Afriyan"},
		}

		tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(executeTransaction)
		wallets.EXPECT().FindByUserIDForUpdate(gomock.Any(), userID).Return(wallet, nil)
		wallets.EXPECT().Update(gomock.Any(), wallet).DoAndReturn(func(_ context.Context, updated *walletdomain.Wallet) error {
			require.EqualValues(t, 60_000, updated.Balance)
			return nil
		})
		codes.EXPECT().TransactionCode().Return("TXN-1")
		codes.EXPECT().ReferenceNumber(userID).Return("REF-1")
		transactions.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, record *transactiondomain.Transaction) error {
			record.ID = transactionID
			require.Equal(t, walletID, record.WalletID)
			require.EqualValues(t, 40_000, record.Amount)
			return nil
		})

		result, err := application.NewService(wallets, transactions, tx, codes).
			Withdraw(context.Background(), userID, 40_000, &remark)
		require.NoError(t, err)
		require.Equal(t, transactionID.String(), result.TransactionID)
		require.Equal(t, "REF-1", result.ReferenceNumber)
		require.EqualValues(t, 60_000, result.Transaction.BalanceNow)
	})

	t.Run("does not update or record when balance is insufficient", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		wallets := mocks.NewMockWalletRepository(ctrl)
		tx := mocks.NewMockTransactionManager(ctrl)
		userID := uuid.New()
		tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(executeTransaction)
		wallets.EXPECT().FindByUserIDForUpdate(gomock.Any(), userID).Return(&walletdomain.Wallet{Balance: 10}, nil)

		_, err := application.NewService(wallets, mocks.NewMockTransactionRepository(ctrl), tx, mocks.NewMockCodeGenerator(ctrl)).
			Withdraw(context.Background(), userID, 20, nil)
		require.ErrorIs(t, err, walletdomain.ErrInsufficientBalance)
	})

	t.Run("returns update error so transaction manager can roll back", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		wallets := mocks.NewMockWalletRepository(ctrl)
		tx := mocks.NewMockTransactionManager(ctrl)
		userID := uuid.New()
		tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(executeTransaction)
		wallet := &walletdomain.Wallet{Balance: 100}
		wallets.EXPECT().FindByUserIDForUpdate(gomock.Any(), userID).Return(wallet, nil)
		wallets.EXPECT().Update(gomock.Any(), wallet).Return(errors.New("update failed"))

		_, err := application.NewService(wallets, mocks.NewMockTransactionRepository(ctrl), tx, mocks.NewMockCodeGenerator(ctrl)).
			Withdraw(context.Background(), userID, 20, nil)
		require.EqualError(t, err, "update failed")
	})
}
