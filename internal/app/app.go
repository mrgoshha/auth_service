package app

import (
	"AuthenticationService/config"
	"AuthenticationService/internal/adapter/dbs/postgres"
	httpApi "AuthenticationService/internal/handler/http"
	"AuthenticationService/pkg/auth"
	emailsender "AuthenticationService/pkg/email_sender"
	"context"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	log             *slog.Logger
	db              *sqlx.DB
	httpServer      *http.Server
	serviceProvider *ServiceProvider
	tokenManager    *auth.Manager
	serviceConfig   *auth.Config
	sender          *emailsender.Sender
}

func NewApp() (*App, error) {
	a := &App{}

	err := a.initDeps()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() {
	go a.runHttpServer()

	a.log.Info("server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.log.Error("failed to stop server", slog.String("error", err.Error()))
	}

	a.log.Info("server stopped")

	if err := a.db.Close(); err != nil {
		a.log.Error("failed to stop storage", slog.String("error", err.Error()))
	}

}

func (a *App) initDeps() error {
	inits := []func() error{
		a.initConfig,
		a.initLogger,
		a.initDb,
		a.initTokenManager,
		a.initEmailSender,
		a.initServiceProvider,
		a.initHttpServer,
	}

	for _, f := range inits {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig() error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	cfg, err := auth.NewConfig()
	if err != nil {
		return err
	}

	a.serviceConfig = cfg
	return nil
}

func (a *App) initLogger() error {
	a.log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return nil
}

func (a *App) initDb() error {
	cfg, err := postgres.NewConfig()
	if err != nil {
		return err
	}

	db, err := postgres.New(cfg)
	if err != nil {
		a.log.Error("failed to init storage", slog.String("error", err.Error()))
		return err
	}

	a.db = db

	return nil
}

func (a *App) initTokenManager() error {
	manager, err := auth.NewManager(a.serviceConfig.SigningKey, a.serviceConfig.AccessTokenTTL)
	if err != nil {
		a.log.Error("failed to init token manager", slog.String("error", err.Error()))
		return err
	}

	a.tokenManager = manager

	return nil
}

func (a *App) initEmailSender() error {
	senderConfig, err := emailsender.NewConfig()
	if err != nil {
		return err
	}

	a.sender = emailsender.NewSender(senderConfig)

	return nil
}

func (a *App) initServiceProvider() error {
	a.serviceProvider = NewServiceProvider(a.log, a.db, a.tokenManager, a.serviceConfig.RefreshTokenTTL, a.sender)
	return nil
}

func (a *App) initHttpServer() error {
	a.serviceProvider.RegisterControllers()

	cfg, err := httpApi.NewConfig()
	if err != nil {
		return err
	}

	srv := httpApi.NewServer(cfg, a.serviceProvider.HttpRouter())

	a.httpServer = srv
	return nil
}

func (a *App) runHttpServer() {
	if err := a.httpServer.ListenAndServe(); err != nil {
		a.log.Error("failed to start server", slog.String("error", err.Error()))
	}
}
