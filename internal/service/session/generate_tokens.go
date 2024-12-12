package session

import (
	"AuthenticationService/internal/model"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *Service) GenerateTokens(id, ip string) (string, string, error) {
	if _, err := s.userService.GetUserById(id); err != nil {
		return "", "", err
	}

	sessionID := generateSessionId()

	accessToken, refreshToken, err := s.generateTokens(id, ip, sessionID)
	if err != nil {
		return "", "", err
	}

	session, err := s.newSession(id, ip, sessionID, refreshToken)
	if err != nil {
		return "", "", err
	}

	if err = s.repository.CreateSession(session); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateSessionId() string {
	return uuid.New().String()
}

func hash(token string) (string, error) {
	hashToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf(`hash token %w`, err)
	}
	return string(hashToken), nil
}

func (s *Service) generateTokens(id, ip, session string) (string, string, error) {
	accessToken, err := s.tokenManager.NewJWT(id, ip, session)
	if err != nil {
		return "", "", fmt.Errorf(`create access token %w`, err)
	}

	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf(`create refresh token %w`, err)
	}

	return accessToken, refreshToken, nil
}

func (s *Service) newSession(id, ip, session, refreshToken string) (*model.Session, error) {
	hashToken, err := hash(refreshToken)
	if err != nil {
		return nil, err
	}
	return &model.Session{
		RefreshToken: hashToken,
		SessionId:    session,
		Ip:           ip,
		UserId:       id,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}, nil
}
