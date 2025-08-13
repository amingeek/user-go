package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"user-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestJWTAuthMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := []byte("testsecret")

	// ساخت توکن معتبر به صورت داینامیک
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"phone": "+12345",
		"exp":   time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	assert.NoError(t, err)

	router := gin.New()
	router.Use(middleware.JWTAuthMiddleware(secret))
	router.GET("/protected", func(c *gin.Context) {
		phone, exists := c.Get("phone")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "phone not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"phone": phone})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "+12345")
}
