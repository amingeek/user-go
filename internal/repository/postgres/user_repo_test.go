package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"user-go/internal/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	// In a real test, you would use a test container
	pool, err := pgxpool.Connect(context.Background(), "postgres://user:pass@localhost:5432/testdb")
	require.NoError(t, err)

	// Clean and migrate
	_, err = pool.Exec(context.Background(), "DROP TABLE IF EXISTS users")
	require.NoError(t, err)

	_, err = pool.Exec(context.Background(), `
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			phone TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP NOT NULL
		)
	`)
	require.NoError(t, err)

	return pool
}

func TestUserRepository(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewUserRepository(pool)

	ctx := context.Background()

	t.Run("CreateUser and GetUserByPhone", func(t *testing.T) {
		phone := "+1234567890"
		user := model.User{
			Phone:     phone,
			CreatedAt: time.Now(),
		}

		createdUser, err := repo.CreateUser(ctx, user)
		require.NoError(t, err)
		assert.NotEmpty(t, createdUser.ID)

		fetchedUser, err := repo.GetUserByPhone(ctx, phone)
		require.NoError(t, err)
		assert.Equal(t, createdUser.ID, fetchedUser.ID)
		assert.Equal(t, phone, fetchedUser.Phone)
	})

	t.Run("GetUsers", func(t *testing.T) {
		// Create multiple users
		phones := []string{"+1111111111", "+2222222222", "+3333333333"}
		for _, phone := range phones {
			_, err := repo.CreateUser(ctx, model.User{
				Phone:     phone,
				CreatedAt: time.Now(),
			})
			require.NoError(t, err)
		}

		users, total, err := repo.GetUsers(ctx, 1, 10, "111")
		require.NoError(t, err)
		assert.Equal(t, 1, len(users))
		assert.Equal(t, "+1111111111", users[0].Phone)
		assert.Equal(t, 1, total)
	})

	t.Run("GetUserByPhone - not found", func(t *testing.T) {
		_, err := repo.GetUserByPhone(ctx, "+9999999999")
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}
