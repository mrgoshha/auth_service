package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"time"
)

//go:generate mockgen -source=manager.go -destination=mocks/mock_manager.go

type TokenManager interface {
	NewJWT(userId, userIP, sessionID string) (string, error)
	Parse(accessToken string) (*TokenPayload, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey     string
	accessTokenTTL time.Duration
}

type tokenClaims struct {
	jwt.StandardClaims
	UserIP    string `json:"ip"`
	SessionId string `json:"session_id"`
}

type TokenPayload struct {
	UserId    string
	UserIP    string
	SessionId string
}

func NewManager(signingKey string, accessTokenTTL time.Duration) (*Manager, error) {
	if signingKey == "" {
		return nil, fmt.Errorf("empty signing key")
	}

	return &Manager{
		signingKey:     signingKey,
		accessTokenTTL: accessTokenTTL,
	}, nil
}

func (m *Manager) NewJWT(userId, userIP, sessionID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.accessTokenTTL).Unix(),
			Subject:   userId,
		},
		UserIP:    userIP,
		SessionId: sessionID,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (*TokenPayload, error) {
	// пропускаем проверку на то что токен истек
	parser := jwt.Parser{SkipClaimsValidation: true}
	token, err := parser.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return nil, fmt.Errorf("parse token")
	}

	return &TokenPayload{
		UserId:    claims.Subject,
		UserIP:    claims.UserIP,
		SessionId: claims.SessionId,
	}, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
