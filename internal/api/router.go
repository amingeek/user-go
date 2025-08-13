package api

import (
	"github.com/amingeek/user-go-complete/internal/jwt"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(h *Handler, jwtManager *jwt.JWTManager) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/otp/request", h.RequestOTP).Methods("POST")
	r.HandleFunc("/otp/login", h.LoginValidateOTP).Methods("POST")

	protectedRouter := r.PathPrefix("/users").Subrouter()
	protectedRouter.Use(func(next http.Handler) http.Handler {
		return AuthMiddleware(jwtManager, next)
	})

	protectedRouter.HandleFunc("/{id}", h.GetUserByID).Methods("GET")
	protectedRouter.HandleFunc("", h.ListUsers).Methods("GET")

	return r
}
