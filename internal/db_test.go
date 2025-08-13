// internal/db_test.go
package internal

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresConnectionAndInsert(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set, skipping DB test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	// درج یک کاربر تستی
	phone := "+989123456789"
	now := time.Now()

	_, err = pool.Exec(ctx, `
		INSERT INTO users (phone, registration_date) 
		VALUES ($1, $2)
		ON CONFLICT (phone) DO NOTHING
	`, phone, now)
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	// خواندن کاربر تستی
	var gotPhone string
	var gotDate time.Time
	err = pool.QueryRow(ctx, `
		SELECT phone, registration_date FROM users WHERE phone = $1
	`, phone).Scan(&gotPhone, &gotDate)
	if err != nil {
		t.Fatalf("failed to query test user: %v", err)
	}

	if gotPhone != phone {
		t.Errorf("expected phone %s, got %s", phone, gotPhone)
	}

	// فقط برای اطمینان
	t.Logf("User %s registered at %v", gotPhone, gotDate)
}
