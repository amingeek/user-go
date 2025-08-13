package user

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUserByPhoneNumber(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)
	phoneNumber := "+989123456789"

	rows := sqlmock.NewRows([]string{"id", "phone_number", "registered_at"}).
		AddRow(1, phoneNumber, time.Now())
	mock.ExpectQuery("SELECT id, phone_number, registered_at FROM users WHERE phone_number = \\$1").
		WithArgs(phoneNumber).
		WillReturnRows(rows)

	user, err := repo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		t.Errorf("expected no error, but got '%s'", err)
	}
	if user == nil || user.PhoneNumber != phoneNumber {
		t.Errorf("expected user with phone number %s, got %v", phoneNumber, user)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRegisterUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)
	newUser := &User{PhoneNumber: "+989121111111"}

	mock.ExpectExec("INSERT INTO users \\(phone_number, registered_at\\) VALUES \\(\\$1, \\$2\\)").
		WithArgs(newUser.PhoneNumber, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.RegisterUser(newUser)
	if err != nil {
		t.Errorf("expected no error, but got '%s'", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
