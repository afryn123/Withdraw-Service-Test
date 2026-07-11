package application

import (
	"context"

	"github.com/afryn123/withdraw-service-test/internal/shared/transaction"
	transactiondomain "github.com/afryn123/withdraw-service-test/internal/transactions/domain"
	walletdomain "github.com/afryn123/withdraw-service-test/internal/wallet/domain"
	"github.com/google/uuid"
)

type WithdrawResult struct {
	TransactionID   string      `json:"transaction_id"`
	ReferenceNumber string      `json:"reference_number"`
	Transaction     Transaction `json:"transaction"`
	User            User        `json:"user"`
}

type Transaction struct {
	Amount     int64   `json:"amount"`
	Type       string  `json:"type"`
	BalanceNow int64   `json:"balance_now"`
	Remark     *string `json:"remark,omitempty"`
}

type User struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type Service struct {
	wallets      walletdomain.Repository
	transactions transactiondomain.Repository
	tx           transaction.Manager
	codes        transactiondomain.CodeGenerator
}

func NewService(wallets walletdomain.Repository, transactions transactiondomain.Repository, tx transaction.Manager, codes transactiondomain.CodeGenerator) *Service {
	return &Service{wallets: wallets, transactions: transactions, tx: tx, codes: codes}
}

func (s *Service) Withdraw(ctx context.Context, userID uuid.UUID, amount int64, remark *string) (WithdrawResult, error) {
	var result WithdrawResult
	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		wallet, err := s.wallets.FindByUserIDForUpdate(txCtx, userID)
		if err != nil {
			return err
		}
		if err := wallet.Withdraw(amount); err != nil {
			return err
		}
		if err := s.wallets.Update(txCtx, wallet); err != nil {
			return err
		}
		record := &transactiondomain.Transaction{
			WalletID: wallet.ID, Amount: amount, Type: transactiondomain.TypeWithdraw,
			TransactionCode: s.codes.TransactionCode(), ReferenceNumber: s.codes.ReferenceNumber(userID),
			Status: transactiondomain.StatusCompleted, Remark: remark,
		}
		if err := s.transactions.Create(txCtx, record); err != nil {
			return err
		}
		result = WithdrawResult{
			TransactionID: record.ID.String(), ReferenceNumber: record.ReferenceNumber,
			Transaction: Transaction{Amount: amount, Type: record.Type, BalanceNow: wallet.Balance, Remark: remark},
			User:        User{UserID: wallet.User.ID.String(), Name: wallet.User.Name},
		}
		return nil
	})
	return result, err
}
