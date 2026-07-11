package application

import (
	"context"

	"github.com/afryn123/withdraw-service-test/internal/shared/transaction"
	userdomain "github.com/afryn123/withdraw-service-test/internal/user/domain"
	walletdomain "github.com/afryn123/withdraw-service-test/internal/wallet/domain"
	"github.com/google/uuid"
)

type CreateUserCommand struct {
	Name      string
	Username  string
	Email     string
	Password  string
	CreatedBy *string
}

type UpdateUserCommand struct {
	UserID   uuid.UUID
	ActorID  uuid.UUID
	Name     *string
	Username *string
	Email    *string
	Password *string
}

type DeleteUserCommand struct {
	UserID  uuid.UUID
	ActorID uuid.UUID
}

type Service struct {
	users   userdomain.Repository
	wallets walletdomain.Repository
	hasher  userdomain.PasswordHasher
	tx      transaction.Manager
}

func NewService(users userdomain.Repository, wallets walletdomain.Repository, hasher userdomain.PasswordHasher, tx transaction.Manager) *Service {
	return &Service{users: users, wallets: wallets, hasher: hasher, tx: tx}
}

func (s *Service) Create(ctx context.Context, command CreateUserCommand) error {
	hashedPassword, err := s.hasher.Hash(command.Password)
	if err != nil {
		return err
	}
	return s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		user := &userdomain.User{
			Name: command.Name, Username: command.Username, Email: command.Email,
			Password: hashedPassword, CreatedBy: command.CreatedBy, UpdatedBy: command.CreatedBy,
		}
		if err := s.users.Create(txCtx, user); err != nil {
			return err
		}
		return s.wallets.Create(txCtx, &walletdomain.Wallet{UserID: user.ID})
	})
}

func (s *Service) Update(ctx context.Context, command UpdateUserCommand) error {
	if command.UserID != command.ActorID {
		return userdomain.ErrForbidden
	}

	user, err := s.users.FindByID(ctx, command.UserID)
	if err != nil {
		return err
	}
	if command.Name != nil {
		user.Name = *command.Name
	}
	if command.Username != nil {
		user.Username = *command.Username
	}
	if command.Email != nil {
		user.Email = *command.Email
	}
	if command.Password != nil {
		hashedPassword, err := s.hasher.Hash(*command.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	updatedBy := command.ActorID.String()
	user.UpdatedBy = &updatedBy
	return s.users.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, command DeleteUserCommand) error {
	if command.UserID != command.ActorID {
		return userdomain.ErrForbidden
	}

	return s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if _, err := s.users.FindByID(txCtx, command.UserID); err != nil {
			return err
		}
		if err := s.wallets.DeleteByUserID(txCtx, command.UserID); err != nil {
			return err
		}
		return s.users.Delete(txCtx, command.UserID)
	})
}
