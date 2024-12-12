package api

import (
	"AuthenticationService/internal/service"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type AuthController struct {
	log     *slog.Logger
	service service.AuthSession
}

func NewAuthController(log *slog.Logger, service service.AuthSession, router *mux.Router) *AuthController {
	ac := &AuthController{
		log:     log,
		service: service,
	}
	router.HandleFunc("/auth/token", ac.GenerateTokens).Methods(http.MethodGet)
	router.HandleFunc("/auth/refresh", ac.Refresh).Methods(http.MethodPost)

	return ac
}
