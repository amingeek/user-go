package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"testing"
	"time"
)

var testRepo *PostgresUserRepository

func setupTestDB(t *testing.T) *pgxpool.Pool {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL env var not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("unable to connect to db: %v", err)
	}
	// پاک کردن جدول و آماده سازی داده تست
	_, err = pool.Exec(context.Background(), "TRUNCATE users")
	if err != nil {
		t.Fatalf("failed to truncate users table: %v", err)
	}
	return pool
}

func TestPostgresUserRepository(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewPostgresUserRepository(pool)

	// Create
	user, err := repo.Create("+1234567890")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if user.Phone != "+1234567890" {
		t.Errorf("expected phone +1234567890, got %s", user.Phone)
	}

	// GetByPhone (exists)
	got, err := repo.GetByPhone("+1234567890")
	if err != nil {
		t.Fatalf("GetByPhone failed: %v", err)
	}
	if got.Phone != "+1234567890" {
		t.Errorf("expected phone +1234567890, got %s", got.Phone)
	}

	// GetByPhone (not exists)
	_, err = repo.GetByPhone("+0000000000")
	if err == nil {
		t.Error("expected error for non-existent user")
	}

	// List
	_, err = repo.Create("+1987654321")
	if err != nil {
		t.Fatalf("Create second user failed: %v", err)
	}
	users, err := repo.List(0, 10, "")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(users) < 2 {
		t.Errorf("expected at least 2 users, got %d", len(users))
	}

	// UpdatePhone
	err = repo.UpdatePhone("+1234567890", "+1111111111")
	if err != nil {
		t.Fatalf("UpdatePhone failed: %v", err)
	}
	updatedUser, err := repo.GetByPhone("+1111111111")
	if err != nil {
		t.Fatalf("GetByPhone after update failed: %v", err)
	}
	if updatedUser.Phone != "+1111111111" {
		t.Errorf("expected phone +1111111111, got %s", updatedUser.Phone)
	}

	// Delete
	err = repo.Delete("+1111111111")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err = repo.GetByPhone("+1111111111")
	if err == nil {
		t.Error("expected error after delete")
	}
}
