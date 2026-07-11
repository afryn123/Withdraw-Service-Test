//go:build integration

package infrastructure_test

import (
	"context"
	"os"
	"testing"

	gormdb "github.com/afryn123/withdraw-service-test/internal/shared/gorm"
	userdomain "github.com/afryn123/withdraw-service-test/internal/user/domain"
	"github.com/afryn123/withdraw-service-test/internal/user/infrastructure"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGormRepositoryLifecycle(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN is required for integration tests")
	}
	db, err := gormdb.Open(dsn)
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&userdomain.User{}))

	repository := infrastructure.NewGormRepository(db)
	ctx := context.Background()
	identity := uuid.NewString()
	user := &userdomain.User{
		Name: "Integration User", Username: "integration-" + identity,
		Email: identity + "@example.com", Password: "hash",
	}
	require.NoError(t, repository.Create(ctx, user))
	t.Cleanup(func() { db.Unscoped().Delete(&userdomain.User{}, "id = ?", user.ID) })

	found, err := repository.FindByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, user.Email, found.Email)

	found.Name = "Updated Integration User"
	require.NoError(t, repository.Update(ctx, found))
	require.NoError(t, repository.Delete(ctx, user.ID))

	_, err = repository.FindByID(ctx, user.ID)
	require.ErrorIs(t, err, userdomain.ErrNotFound)
}
