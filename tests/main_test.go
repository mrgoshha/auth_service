package tests

import (
	"AuthenticationService/internal/app"
	"AuthenticationService/pkg/auth"
	sendermocks "AuthenticationService/pkg/email_sender/mocks"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"os"
	"testing"
	"time"
)

type APITestSuite struct {
	suite.Suite

	db *sqlx.DB

	tokenManager *auth.Manager

	sender *sendermocks.EmailSender

	serviceProvider *app.ServiceProvider
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	dataSource := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		"test", "test", "localhost", "5430", "authServiceTest")

	if db, err := sqlx.Connect("postgres", dataSource); err != nil {
		s.FailNow("Failed to connect to postgres", err)
	} else {
		s.db = db
	}

	s.initDeps()

	if err := s.initDB(); err != nil {
		s.FailNow("Failed to create and populate DB", err)
	}
}

func (s *APITestSuite) TearDownSuite() {
	s.db.Close()
}

func (s *APITestSuite) initDeps() {
	// Init domain deps
	d, _ := time.ParseDuration("2h")
	if tokenManager, err := auth.NewManager("signing_key", d); err != nil {
		s.FailNow("Failed to initialize token manager", err)
	} else {
		s.tokenManager = tokenManager
	}

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	s.sender = new(sendermocks.EmailSender)

	s.serviceProvider = app.NewServiceProvider(log, s.db, s.tokenManager, d, s.sender)

	s.serviceProvider.RegisterControllers()

}

func (s *APITestSuite) initDB() error {
	createUsers := `	CREATE TABLE IF NOT EXISTS users(
						user_id VARCHAR PRIMARY KEY,
						email VARCHAR NOT NULL)`

	_, err := s.db.Exec(createUsers)
	if err != nil {
		s.FailNow("Failed to create users table", err)
	}

	insertUser := `INSERT INTO users VALUES ('92503be1-3a71-4135-a7ea-b42361957c56', 'test@mail.ru');`

	_, err = s.db.Exec(insertUser)
	if err != nil {
		s.FailNow("Failed to insert user", err)
	}

	createSessions := `	CREATE TABLE IF NOT EXISTS sessions (
						id SERIAL PRIMARY KEY,
						session_id VARCHAR NOT NULL,
						refresh_token VARCHAR NOT NULL,
						ip VARCHAR NOT NULL,
						user_id VARCHAR REFERENCES users ON DELETE CASCADE,
						expires_at TIMESTAMP NOT NULL)`

	_, err = s.db.Exec(createSessions)
	if err != nil {
		s.FailNow("Failed to create sessions table", err)
	}

	expiresAt := time.Now().Add(time.Hour)
	insertSession := `	INSERT INTO sessions (session_id, refresh_token, ip, user_id, expires_at)
						VALUES ($1, $2, $3, $4, $5)`

	_, err = s.db.Exec(insertSession, "5b0f8c90-764b-4971-b82a-48f9f4c71705",
		"$2a$10$t3.JHeLr8Rp8.1jbB7VAG.Ne90Bxoq2vdbAkO8n1uME/qXtQkzJV6",
		"ip",
		"92503be1-3a71-4135-a7ea-b42361957c56",
		expiresAt)
	if err != nil {
		s.FailNow("Failed to insert user", err)
	}

	return nil
}
