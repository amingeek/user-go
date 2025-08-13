package repository

import (
	"testing"
)

func TestInMemoryUserRepository_CreateAndGet(t *testing.T) {
	repo := NewInMemoryUserRepository()

	phone := "+1234567890"

	user, err := repo.Create(phone)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Phone != phone {
		t.Errorf("expected phone %s, got %s", phone, user.Phone)
	}

	gotUser, err := repo.GetByPhone(phone)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotUser.Phone != phone {
		t.Errorf("expected phone %s, got %s", phone, gotUser.Phone)
	}
}

func TestInMemoryUserRepository_List(t *testing.T) {
	repo := NewInMemoryUserRepository()

	phones := []string{"+111", "+222", "+333"}
	for _, p := range phones {
		_, err := repo.Create(p)
		if err != nil {
			t.Fatalf("unexpected error creating user: %v", err)
		}
	}

	users, err := repo.List(0, 10, "+2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
	if users[0].Phone != "+222" {
		t.Errorf("expected phone +222, got %s", users[0].Phone)
	}
}

func TestInMemoryUserRepository_UpdatePhone(t *testing.T) {
	repo := NewInMemoryUserRepository()

	_, _ = repo.Create("+111")
	_, _ = repo.Create("+222")

	// تغییر شماره معتبر
	err := repo.UpdatePhone("+111", "+333")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = repo.GetByPhone("+111")
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}

	user, err := repo.GetByPhone("+333")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Phone != "+333" {
		t.Errorf("expected phone +333, got %s", user.Phone)
	}

	// تغییر به شماره تکراری
	err = repo.UpdatePhone("+222", "+333")
	if err == nil {
		t.Errorf("expected error for duplicate new phone")
	}

	// تغییر شماره غیر موجود
	err = repo.UpdatePhone("+999", "+444")
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestInMemoryUserRepository_Delete(t *testing.T) {
	repo := NewInMemoryUserRepository()

	_, _ = repo.Create("+111")

	err := repo.Delete("+111")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = repo.GetByPhone("+111")
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}

	// حذف شماره غیر موجود
	err = repo.Delete("+222")
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
