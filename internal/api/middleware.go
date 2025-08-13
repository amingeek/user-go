package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/amingeek/user-go-complete/internal/jwt"
)

type contextKey string

const userClaimsKey contextKey = "user_claims"

func AuthMiddleware(jwtManager *jwt.JWTManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Malformed token", http.StatusUnauthorized)
			return
		}

		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
