package domain_test

import (
	"testing"

	walletdomain "github.com/afryn123/withdraw-service-test/internal/wallet/domain"
	"github.com/stretchr/testify/require"
)

func TestWalletWithdraw(t *testing.T) {
	t.Run("deducts balance", func(t *testing.T) {
		wallet := walletdomain.Wallet{Balance: 100}
		require.NoError(t, wallet.Withdraw(40))
		require.EqualValues(t, 60, wallet.Balance)
	})

	t.Run("rejects insufficient balance without mutating it", func(t *testing.T) {
		wallet := walletdomain.Wallet{Balance: 10}
		err := wallet.Withdraw(20)
		require.ErrorIs(t, err, walletdomain.ErrInsufficientBalance)
		require.EqualValues(t, 10, wallet.Balance)
	})

	t.Run("rejects non-positive amount", func(t *testing.T) {
		wallet := walletdomain.Wallet{Balance: 100}
		err := wallet.Withdraw(0)
		require.ErrorIs(t, err, walletdomain.ErrInvalidWithdrawAmount)
		require.EqualValues(t, 100, wallet.Balance)
	})
}
