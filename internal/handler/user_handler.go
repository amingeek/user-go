package handler

import (
	"net/http"
	"strconv"
	"user-go/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	phoneVal, exists := c.Get("phone")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "phone not found in context"})
		return
	}

	phone, ok := phoneVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid phone type"})
		return
	}

	user, err := h.userRepo.GetByPhone(phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	search := c.Query("search")

	offset := 0
	limit := 10

	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	users, err := h.userRepo.List(offset, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
