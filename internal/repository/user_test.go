package repository_test

import (
	"testing"
	"user-go/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryUserRepository_CreateAndGet(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()

	phone := "+1234567890"

	// ایجاد کاربر جدید
	user, err := repo.Create(phone)
	require.NoError(t, err)
	assert.Equal(t, phone, user.Phone)

	// تلاش دوباره برای ایجاد همان کاربر باید خطا بده
	_, err = repo.Create(phone)
	assert.Error(t, err)

	// گرفتن کاربر
	got, err := repo.GetByPhone(phone)
	require.NoError(t, err)
	assert.Equal(t, user.Phone, got.Phone)

	// گرفتن کاربر ناموجود باید خطا بده
	_, err = repo.GetByPhone("+0000000000")
	assert.Equal(t, repository.ErrUserNotFound, err)
}

func TestInMemoryUserRepository_List(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()

	// ایجاد چند کاربر
	phones := []string{"+111", "+222", "+333"}
	for _, p := range phones {
		_, _ = repo.Create(p)
	}

	// لیست همه کاربران بدون فیلتر
	users, err := repo.List(0, 10, "")
	require.NoError(t, err)
	assert.Len(t, users, 3)

	// لیست با جستجو
	users, err = repo.List(0, 10, "+2")
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "+222", users[0].Phone)

	// لیست با صفحه‌بندی
	users, err = repo.List(1, 1, "")
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "+222", users[0].Phone)
}
