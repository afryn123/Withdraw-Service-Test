package infrastructure_test

import (
	"testing"

	"github.com/afryn123/withdraw-service-test/internal/user/infrastructure"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptPassword(t *testing.T) {
	passwords := infrastructure.NewBcryptPassword(bcrypt.MinCost)
	hash, err := passwords.Hash("secret12")
	require.NoError(t, err)
	require.NotEqual(t, "secret12", hash)
	require.True(t, passwords.Verify("secret12", hash))
	require.False(t, passwords.Verify("wrong-password", hash))
}
