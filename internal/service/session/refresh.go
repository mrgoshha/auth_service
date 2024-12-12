package session

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	message = "Замечена подозрительная активность в вашем аккаунте.\n\nВозможно, кто-то вошёл в ваш аккаунт с другого IP-адреса. " +
		"\n\nПроверьте безопасность аккаунта и смените пароль, если это необходимо. " +
		"Если вы сами авторизовались с нового устройства или места, то всё в порядке."
)

// Refresh парсим access токен, чтобы получить id сессии
// С помощью id сессии ищем refresh токен в бд, таким образом проверяем что токены парные
// далее сравниваем refresh токены и если они совпадают проверяем время жизни
// проверяем ip и если он не совпадает отправляем warning email
func (s *Service) Refresh(aToken, rToken, ip string) (string, string, error) {
	tokenPayload, err := s.tokenManager.Parse(aToken)
	if err != nil {
		return "", "", err
	}

	user, err := s.userService.GetUserById(tokenPayload.UserId)
	if err != nil {
		return "", "", err
	}

	session, err := s.repository.GetSessionBySessionId(tokenPayload.SessionId)
	if err != nil {
		return "", "", fmt.Errorf(`unpaired tokens %w`, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(session.RefreshToken), []byte(rToken)); err != nil {
		return "", "", fmt.Errorf(`token is not valid %w`, err)
	}

	if time.Now().After(session.ExpiresAt) {
		return "", "", fmt.Errorf(`token is expired`)
	}

	if session.Ip != ip {
		if err = s.sender.SendEmail(user.Email, message); err != nil {
			return "", "", err
		}
	}

	newSessionId := generateSessionId()

	accessToken, refreshToken, err := s.generateTokens(tokenPayload.UserId, ip, newSessionId)
	if err != nil {
		return "", "", err
	}

	hashRefreshToken, err := hash(refreshToken)
	if err != nil {
		return "", "", err
	}

	session.SessionId = newSessionId
	session.RefreshToken = hashRefreshToken
	session.Ip = ip
	session.ExpiresAt = time.Now().Add(s.refreshTokenTTL)

	if err = s.repository.UpdateSession(session); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
