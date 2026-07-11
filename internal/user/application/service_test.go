package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/afryn123/withdraw-service-test/internal/mocks"
	"github.com/afryn123/withdraw-service-test/internal/user/application"
	userdomain "github.com/afryn123/withdraw-service-test/internal/user/domain"
	walletdomain "github.com/afryn123/withdraw-service-test/internal/wallet/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func executeUserTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestServiceCreate(t *testing.T) {
	t.Run("creates user and wallet in one transaction", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		users := mocks.NewMockUserRepository(ctrl)
		wallets := mocks.NewMockWalletRepository(ctrl)
		hasher := mocks.NewMockPasswordHasher(ctrl)
		tx := mocks.NewMockTransactionManager(ctrl)
		userID := uuid.New()

		hasher.EXPECT().Hash("secret12").Return("hashed", nil)
		tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(context.Context) error) error { return fn(ctx) },
		)
		users.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, user *userdomain.User) error {
			user.ID = userID
			require.Equal(t, "hashed", user.Password)
			return nil
		})
		wallets.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, wallet *walletdomain.Wallet) error {
			require.Equal(t, userID, wallet.UserID)
			require.Zero(t, wallet.Balance)
			return nil
		})

		service := application.NewService(users, wallets, hasher, tx)
		err := service.Create(context.Background(), application.CreateUserCommand{
			Name: "Afriyan", Username: "afryn", Email: "a@example.com", Password: "secret12",
		})
		require.NoError(t, err)
	})

	t.Run("stops before transaction when password hashing fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		hasher := mocks.NewMockPasswordHasher(ctrl)
		hasher.EXPECT().Hash("secret12").Return("", errors.New("hash failed"))
		service := application.NewService(
			mocks.NewMockUserRepository(ctrl), mocks.NewMockWalletRepository(ctrl), hasher, mocks.NewMockTransactionManager(ctrl),
		)
		err := service.Create(context.Background(), application.CreateUserCommand{Password: "secret12"})
		require.EqualError(t, err, "hash failed")
	})
}

func TestServiceUpdate(t *testing.T) {
	t.Run("updates provided fields and hashes a new password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		users := mocks.NewMockUserRepository(ctrl)
		hasher := mocks.NewMockPasswordHasher(ctrl)
		userID := uuid.New()
		name, password := "Afriyan Updated", "new-secret"
		user := &userdomain.User{
			ID: userID, Name: "Afriyan", Username: "afryn",
			Email: "a@example.com", Password: "old-hash",
		}

		users.EXPECT().FindByID(gomock.Any(), userID).Return(user, nil)
		hasher.EXPECT().Hash(password).Return("new-hash", nil)
		users.EXPECT().Update(gomock.Any(), user).DoAndReturn(func(_ context.Context, updated *userdomain.User) error {
			require.Equal(t, name, updated.Name)
			require.Equal(t, "afryn", updated.Username)
			require.Equal(t, "a@example.com", updated.Email)
			require.Equal(t, "new-hash", updated.Password)
			require.NotNil(t, updated.UpdatedBy)
			require.Equal(t, userID.String(), *updated.UpdatedBy)
			return nil
		})

		service := application.NewService(
			users, mocks.NewMockWalletRepository(ctrl), hasher, mocks.NewMockTransactionManager(ctrl),
		)
		err := service.Update(context.Background(), application.UpdateUserCommand{
			UserID: userID, ActorID: userID, Name: &name, Password: &password,
		})
		require.NoError(t, err)
	})

	t.Run("rejects an actor editing another user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := application.NewService(
			mocks.NewMockUserRepository(ctrl), mocks.NewMockWalletRepository(ctrl),
			mocks.NewMockPasswordHasher(ctrl), mocks.NewMockTransactionManager(ctrl),
		)

		err := service.Update(context.Background(), application.UpdateUserCommand{
			UserID: uuid.New(), ActorID: uuid.New(),
		})
		require.ErrorIs(t, err, userdomain.ErrForbidden)
	})

	t.Run("returns user lookup error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		users := mocks.NewMockUserRepository(ctrl)
		userID := uuid.New()
		users.EXPECT().FindByID(gomock.Any(), userID).Return(nil, userdomain.ErrNotFound)
		service := application.NewService(
			users, mocks.NewMockWalletRepository(ctrl), mocks.NewMockPasswordHasher(ctrl), mocks.NewMockTransactionManager(ctrl),
		)

		err := service.Update(context.Background(), application.UpdateUserCommand{UserID: userID, ActorID: userID})
		require.ErrorIs(t, err, userdomain.ErrNotFound)
	})

	t.Run("does not update when password hashing fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		users := mocks.NewMockUserRepository(ctrl)
		hasher := mocks.NewMockPasswordHasher(ctrl)
		userID := uuid.New()
		password := "new-secret"
		users.EXPECT().FindByID(gomock.Any(), userID).Return(&userdomain.User{ID: userID}, nil)
		hasher.EXPECT().Hash(password).Return("", errors.New("hash failed"))
		service := application.NewService(
			users, mocks.NewMockWalletRepository(ctrl), hasher, mocks.NewMockTransactionManager(ctrl),
		)

		err := service.Update(context.Background(), application.UpdateUserCommand{
			UserID: userID, ActorID: userID, Password: &password,
		})
		require.EqualError(t, err, "hash failed")
	})
}

func TestServiceDelete(t *testing.T) {
	t.Run("soft deletes wallet and user in one transaction", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		users := mocks.NewMockUserRepository(ctrl)
		wallets := mocks.NewMockWalletRepository(ctrl)
		tx := mocks.NewMockTransactionManager(ctrl)
		userID := uuid.New()

		tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(executeUserTransaction)
		gomock.InOrder(
			users.EXPECT().FindByID(gomock.Any(), userID).Return(&userdomain.User{ID: userID}, nil),
			wallets.EXPECT().DeleteByUserID(gomock.Any(), userID).Return(nil),
			users.EXPECT().Delete(gomock.Any(), userID).Return(nil),
		)

		service := application.NewService(users, wallets, mocks.NewMockPasswordHasher(ctrl), tx)
		err := service.Delete(context.Background(), application.DeleteUserCommand{UserID: userID, ActorID: userID})
		require.NoError(t, err)
	})

	t.Run("rejects an actor deleting another user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := application.NewService(
			mocks.NewMockUserRepository(ctrl), mocks.NewMockWalletRepository(ctrl),
			mocks.NewMockPasswordHasher(ctrl), mocks.NewMockTransactionManager(ctrl),
		)

		err := service.Delete(context.Background(), application.DeleteUserCommand{
			UserID: uuid.New(), ActorID: uuid.New(),
		})
		require.ErrorIs(t, err, userdomain.ErrForbidden)
	})

	t.Run("stops when user does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		users := mocks.NewMockUserRepository(ctrl)
		tx := mocks.NewMockTransactionManager(ctrl)
		userID := uuid.New()
		tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(executeUserTransaction)
		users.EXPECT().FindByID(gomock.Any(), userID).Return(nil, userdomain.ErrNotFound)

		service := application.NewService(users, mocks.NewMockWalletRepository(ctrl), mocks.NewMockPasswordHasher(ctrl), tx)
		err := service.Delete(context.Background(), application.DeleteUserCommand{UserID: userID, ActorID: userID})
		require.ErrorIs(t, err, userdomain.ErrNotFound)
	})

	t.Run("returns wallet error so the transaction rolls back", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		users := mocks.NewMockUserRepository(ctrl)
		wallets := mocks.NewMockWalletRepository(ctrl)
		tx := mocks.NewMockTransactionManager(ctrl)
		userID := uuid.New()
		tx.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(executeUserTransaction)
		users.EXPECT().FindByID(gomock.Any(), userID).Return(&userdomain.User{ID: userID}, nil)
		wallets.EXPECT().DeleteByUserID(gomock.Any(), userID).Return(errors.New("delete wallet failed"))

		service := application.NewService(users, wallets, mocks.NewMockPasswordHasher(ctrl), tx)
		err := service.Delete(context.Background(), application.DeleteUserCommand{UserID: userID, ActorID: userID})
		require.EqualError(t, err, "delete wallet failed")
	})
}
