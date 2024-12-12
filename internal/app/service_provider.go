package app

import (
	sessionRepo "AuthenticationService/internal/adapter/dbs/postgres/session"
	userRepo "AuthenticationService/internal/adapter/dbs/postgres/user"
	"AuthenticationService/internal/handler/http"
	"AuthenticationService/internal/handler/http/api"
	"AuthenticationService/internal/service/session"
	"AuthenticationService/internal/service/user"
	"AuthenticationService/pkg/auth"
	emailsender "AuthenticationService/pkg/email_sender"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"time"
)

type ServiceProvider struct {
	authController    *api.AuthController
	authService       *session.Service
	userService       *user.Service
	sessionRepository *sessionRepo.Repository
	userRepository    *userRepo.Repository
	httpRouter        *mux.Router
	log               *slog.Logger
	db                *sqlx.DB
	tokenManager      *auth.Manager
	refreshTokenTTL   time.Duration
	sender            emailsender.EmailSender
}

func NewServiceProvider(log *slog.Logger, db *sqlx.DB, t *auth.Manager, refreshTokenTTL time.Duration, sender emailsender.EmailSender) *ServiceProvider {
	return &ServiceProvider{
		log:             log,
		db:              db,
		tokenManager:    t,
		refreshTokenTTL: refreshTokenTTL,
		sender:          sender,
	}
}

func (s *ServiceProvider) AuthController() *api.AuthController {
	if s.authController == nil {
		s.authController = api.NewAuthController(
			s.log, s.AuthService(), s.HttpRouter(),
		)
	}
	return s.authController
}

func (s *ServiceProvider) AuthService() *session.Service {
	if s.authService == nil {
		s.authService = session.NewService(
			s.SessionRepository(), s.tokenManager, s.refreshTokenTTL, s.UserService(), s.sender,
		)
	}
	return s.authService
}

func (s *ServiceProvider) UserService() *user.Service {
	if s.userService == nil {
		s.userService = user.NewService(
			s.UserRepository(),
		)
	}
	return s.userService
}

func (s *ServiceProvider) SessionRepository() *sessionRepo.Repository {
	if s.sessionRepository == nil {
		s.sessionRepository = sessionRepo.NewRepository(
			s.db,
		)
	}
	return s.sessionRepository
}

func (s *ServiceProvider) UserRepository() *userRepo.Repository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewRepository(
			s.db,
		)
	}
	return s.userRepository
}

func (s *ServiceProvider) HttpRouter() *mux.Router {
	if s.httpRouter == nil {
		s.httpRouter = http.NewRouter(s.log)
	}
	return s.httpRouter
}

func (s *ServiceProvider) RegisterControllers() {
	if s.authController == nil {
		s.AuthController()
	}
}
