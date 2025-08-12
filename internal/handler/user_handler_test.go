package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-go/internal/handler"
	"user-go/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupUserRouter(repo repository.UserRepository) *gin.Engine {
	r := gin.Default()
	userHandler := handler.NewUserHandler(repo)
	r.GET("/users", userHandler.ListUsers)
	return r
}

func TestListUsersHandler(t *testing.T) {
	// ایجاد ریپوزیتوری و اضافه کردن چند کاربر
	repo := repository.NewInMemoryUserRepository()
	_, _ = repo.Create("+111")
	_, _ = repo.Create("+222")
	_, _ = repo.Create("+333")

	r := setupUserRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/users?offset=1&limit=2", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var users []repository.User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// انتظار داریم کاربران با ایندکس 1 و 2 رو دریافت کنیم (+222 و +333)
	phones := []string{}
	for _, u := range users {
		phones = append(phones, u.Phone)
	}

	assert.Contains(t, phones, "+222")
	assert.Contains(t, phones, "+333")
	assert.NotContains(t, phones, "+111")
	assert.Len(t, users, 2)
}
