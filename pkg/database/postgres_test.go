package database

import (
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewPostgresDB(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("skipping integration test, TEST_DATABASE_URL not set")
	}

	dbConn, err := NewPostgresDB(os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if dbConn == nil {
		t.Error("expected a database connection, got nil")
	}
}

func TestCreateUserTable(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = CreateUserTable(db)
	if err != nil {
		t.Errorf("expected no error, got '%s'", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
