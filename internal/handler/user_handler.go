package handler

import (
	"net/http"
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
	phone, exists := c.Get("phone")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userRepo.GetByPhone(phone.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	phone := c.Param("phone")
	user, err := h.userRepo.GetByPhone(phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	offset := 0
	limit := 10
	search := c.Query("search")

	users, err := h.userRepo.List(offset, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) EditUser(c *gin.Context) {
	phone := c.Param("phone")

	var req struct {
		NewPhone string `json:"new_phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "new_phone is required"})
		return
	}

	err := h.userRepo.UpdatePhone(phone, req.NewPhone)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	phone := c.Param("phone")

	err := h.userRepo.Delete(phone)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
