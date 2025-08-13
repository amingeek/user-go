package middleware

import (
	"net/http"
	"strconv"
	"time"

	"user-go/internal/service"

	"github.com/go-chi/chi/v5/middleware"
)

type RateLimitMiddleware struct {
	rateLimitService service.RateLimitService
}

func NewRateLimitMiddleware(rateLimitService service.RateLimitService) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		rateLimitService: rateLimitService,
	}
}

func (m *RateLimitMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client IP address
		clientIP := middleware.GetRealIP(r)

		// Check rate limit
		limited, err := m.rateLimitService.Check(r.Context(), "ip:"+clientIP)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if limited {
			w.Header().Set("Retry-After", strconv.Itoa(int(10*time.Minute.Seconds())))
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
