package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-go/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() (*gin.Engine, *repository.InMemoryUserRepository) {
	r := gin.Default()
	userRepo := repository.NewInMemoryUserRepository()
	userHandler := NewUserHandler(userRepo)

	r.GET("/users/:phone", userHandler.GetUser)
	r.GET("/users", userHandler.ListUsers)
	r.PUT("/users/:phone", userHandler.EditUser)
	r.DELETE("/users/:phone", userHandler.DeleteUser)

	return r, userRepo
}

func TestGetUser(t *testing.T) {
	r, repo := setupRouter()

	_, _ = repo.Create("+123")

	req, _ := http.NewRequest("GET", "/users/+123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestListUsers(t *testing.T) {
	r, repo := setupRouter()

	_, _ = repo.Create("+111")
	_, _ = repo.Create("+222")

	req, _ := http.NewRequest("GET", "/users?search=+2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var users []repository.User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "+222", users[0].Phone)
}

func TestEditUser(t *testing.T) {
	r, repo := setupRouter()

	_, _ = repo.Create("+111")

	body := map[string]string{"new_phone": "+999"}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", "/users/+111", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	user, err := repo.GetByPhone("+999")
	assert.NoError(t, err)
	assert.Equal(t, "+999", user.Phone)
}

func TestDeleteUser(t *testing.T) {
	r, repo := setupRouter()

	_, _ = repo.Create("+111")

	req, _ := http.NewRequest("DELETE", "/users/+111", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	_, err := repo.GetByPhone("+111")
	assert.ErrorIs(t, err, repository.ErrUserNotFound)
}
