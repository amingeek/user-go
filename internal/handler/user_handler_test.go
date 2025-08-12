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

func setupRouterWithUserHandler(userRepo repository.UserRepository) *gin.Engine {
	r := gin.Default()
	userHandler := handler.NewUserHandler(userRepo)

	r.GET("/users/:phone", userHandler.GetUser)
	r.GET("/users", userHandler.ListUsers)

	return r
}

func TestGetUser_Success(t *testing.T) {
	userRepo := repository.NewInMemoryUserRepository()
	userRepo.Create("+1234567890")

	router := setupRouterWithUserHandler(userRepo)

	req, _ := http.NewRequest("GET", "/users/+1234567890", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var data map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &data)
	assert.NoError(t, err)
	assert.Equal(t, "+1234567890", data["phone"])
}

func TestGetUser_NotFound(t *testing.T) {
	userRepo := repository.NewInMemoryUserRepository()

	router := setupRouterWithUserHandler(userRepo)

	req, _ := http.NewRequest("GET", "/users/+0000000000", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)

	var data map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &data)
	assert.NoError(t, err)
	assert.Contains(t, data["error"], "not found")
}

func TestListUsers(t *testing.T) {
	userRepo := repository.NewInMemoryUserRepository()
	userRepo.Create("+111")
	userRepo.Create("+222")
	userRepo.Create("+333")

	router := setupRouterWithUserHandler(userRepo)

	req, _ := http.NewRequest("GET", "/users?offset=0&limit=2&search=+2", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var data struct {
		Users []struct {
			Phone            string `json:"phone"`
			RegistrationDate string `json:"registration_date"`
		} `json:"users"`
	}

	err := json.Unmarshal(resp.Body.Bytes(), &data)
	assert.NoError(t, err)
	assert.Len(t, data.Users, 1)
	assert.Equal(t, "+222", data.Users[0].Phone)
}
